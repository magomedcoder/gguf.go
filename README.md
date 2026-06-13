# gguf.go - запуск ML-моделей в формате **GGUF** на чистом **Go**.

> **Ранний этап разработки.**

**gguf.go** - лёгковесный способ запуска GGUF-моделей на языке Go без использования llama.cpp.

Формат **GGUF** используется в экосистеме llama.cpp. 

---

## Что уже работает

- парсинг GGUF v2/v3 (`info`, `inspect`);
- деквантизация Q8_0, загрузка весов (`quant`, `tensor`, `weights`);
- базовые ops: RoPE, RMSNorm, GQA attention, SwiGLU;
- forward pass Qwen3 + KV-cache (`model/qwen3`, `runtime`);
- tokenizer BPE из метаданных GGUF (`tokenizer`);
- генерация текста: `gguf run` (prefill + greedy/temperature/top-k/top-p);
- Qwen3 Instruct: `--chat` для chat template.

---

На текущем этапе для разработки и тестирования используется Qwen3-0.6B-Q8_0.gguf

```bash
mkdir -p models

curl -L -o models/Qwen3-0.6B-Q8_0.gguf https://huggingface.co/Qwen/Qwen3-0.6B-GGUF/resolve/main/Qwen3-0.6B-Q8_0.gguf
```

```bash
go build -o build/gguf ./cmd/gguf
```

### `gguf info`

Краткая сводка о модели: версия GGUF, архитектура, имя, число тензоров, размер весов, длина контекста.

```bash
./build/gguf info -m ./models/Qwen3-0.6B-Q8_0.gguf
```

| Флаг | Описание |
|------|----------|
| `-m` | путь к файлу `.gguf` |

### `gguf inspect`

Полный дамп метаданных и списка тензоров (имя, тип, размерности, размер в байтах).

```bash
./build/gguf inspect ./models/Qwen3-0.6B-Q8_0.gguf
```

Аргумент — путь к файлу, без флагов.

### `gguf run`

Генерация текста: prefill промпта -> autoregressive decode -> вывод в stdout.

```bash
./build/gguf run -m ./models/Qwen3-0.6B-Q8_0.gguf --chat -p "Привет" -n 64
```

| Флаг | По умолчанию | Описание |
|------|--------------|----------|
| `-m` | — | путь к файлу `.gguf` |
| `-p` | — | текст промпта |
| `-n` | `128` | максимум новых токенов |
| `--temp` | `0` | температура sampling (`0` = greedy) |
| `--top-k` | `0` | top-k (`0` = выключено) |
| `--top-p` | `1` | nucleus sampling (`1` = выключено) |
| `--seed` | `0` | seed PRNG |
| `--chat` | `false` | обернуть промпт в Qwen chat template |

Для **Qwen3 Instruct** используйте `--chat`, иначе модель ответит некорректно.

Пример с sampling:

```bash
./build/gguf run -m ./models/Qwen3-0.6B-Q8_0.gguf --chat -p "Привет" -n 64 --temp 0.7 --top-k 40 --top-p 0.9 --seed 42
```

---

### Утилиты для отладки

#### `debugtok`

Проверяет encode промпта и logits после prefill: top-5 токенов и greedy-следующий.

```bash
go run ./cmd/debugtok ./models/Qwen3-0.6B-Q8_0.gguf "Hello"
go run ./cmd/debugtok ./models/Qwen3-0.6B-Q8_0.gguf '<|im_start|>user
Hello
<|im_start|>assistant
'
```

#### `vocab`

Показывает конфиг Qwen3 (`head_dim`, число heads) и ID special tokens в словаре (`<|im_start|>`, `<|endoftext|>`, etc.).

```bash
go run ./cmd/vocab ./models/Qwen3-0.6B-Q8_0.gguf
```


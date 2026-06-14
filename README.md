# gguf.go - запуск ML-моделей в формате **GGUF** на чистом **Go**.

> **Ранний этап разработки.**

**gguf.go** - лёгковесный способ запуска GGUF-моделей на языке Go без использования llama.cpp.

Формат **GGUF** используется в экосистеме llama.cpp.

---

## Что уже работает

- парсинг GGUF v2/v3 (`info`, `inspect`);
- деквантизация Q8_0, загрузка весов (`quant`, `tensor`, `weights`);
- базовые ops: RoPE, RMSNorm, GQA attention, SwiGLU;
- forward pass Qwen3 + KV-cache (`pkg/model/qwen3`, `pkg/runtime`);
- tokenizer BPE из метаданных GGUF (`pkg/tokenizer`);
- генерация текста: `gguf run` (prefill + greedy/temperature/top-k/top-p);
- Qwen3 Instruct: `--chat` для chat template.

---

На текущем этапе для разработки и тестирования используется Qwen3-0.6B-Q8_0.gguf

```bash
mkdir -p models

curl -L -o models/Qwen3-0.6B-Q8_0.gguf https://huggingface.co/Qwen/Qwen3-0.6B-GGUF/resolve/main/Qwen3-0.6B-Q8_0.gguf
```

---

### Сборка локально (Go 1.26)

```bash
go build -o build/gguf ./cmd/gguf
```

### Сборка через Docker (Linux, macOS, Windows)

Кросс-компиляция всех платформ в одном образе:

```bash
docker build -t gguf-build .

docker run --rm -v "$(pwd)/build:/out" gguf-build
```

Результат в `build/`:

| Платформа     | Файл                           |
|---------------|--------------------------------|
| Linux amd64   | `build/linux-amd64/gguf`       |
| Linux arm64   | `build/linux-arm64/gguf`       |
| Windows amd64 | `build/windows-amd64/gguf.exe` |
| Windows arm64 | `build/windows-arm64/gguf.exe` |
| macOS amd64   | `build/darwin-amd64/gguf`      |
| macOS arm64   | `build/darwin-arm64/gguf`      |

> **Примечание.** Путь к бинарнику зависит от способа сборки:
> - локально: `./build/gguf`
> - через Docker укажите платформу - `./build/<os>-<arch>/gguf` (на Linux amd64: `./build/linux-amd64/gguf`, на macOS arm64: `./build/darwin-arm64/gguf`, на Windows: `./build/windows-amd64/gguf.exe`)

---

### `gguf info`

Краткая сводка о модели: версия GGUF, архитектура, имя, число тензоров, размер весов, длина контекста.

```bash
./build/gguf info -m ./models/Qwen3-0.6B-Q8_0.gguf
```

| Флаг | Описание             |
|------|----------------------|
| `-m` | путь к файлу `.gguf` |

### `gguf inspect`

Полный дамп метаданных и списка тензоров (имя, тип, размерности, размер в байтах).

```bash
./build/gguf inspect ./models/Qwen3-0.6B-Q8_0.gguf
```

Аргумент - путь к файлу, без флагов.

### `gguf run`

Генерация текста: prefill промпта -> autoregressive decode -> вывод в stdout.

```bash
./build/gguf run -m ./models/Qwen3-0.6B-Q8_0.gguf --chat -p "Привет" -n 64
```

| Флаг      | По умолчанию | Описание                             |
|-----------|--------------|--------------------------------------|
| `-m`      | -            | путь к файлу `.gguf`                 |
| `-p`      | -            | текст промпта                        |
| `-n`      | `128`        | максимум новых токенов               |
| `--temp`  | `0`          | температура sampling (`0` = greedy)  |
| `--top-k` | `0`          | top-k (`0` = выключено)              |
| `--top-p` | `1`          | nucleus sampling (`1` = выключено)   |
| `--seed`  | `0`          | seed PRNG                            |
| `--chat`  | `false`      | обернуть промпт в Qwen chat template |

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

---

### Использование как библиотеки

```go
import "github.com/magomedcoder/gguf.go"

// Парсинг файла
r, err := gguf.OpenFile("./models/Qwen3-0.6B-Q8_0.gguf")

// Inference
engine, err := gguf.Load("./models/Qwen3-0.6B-Q8_0.gguf")
ctx, err := engine.NewContext()
text, err := ctx.Generate("Привет", gguf.GenerateParams{
    MaxTokens: 128,
    Sampler:   gguf.Greedy,
})
```

Реализация разбита по пакетам в pkg/ (format, quant, ops, model, runtime и тд).
Их можно импортировать напрямую при расширении или отладке.
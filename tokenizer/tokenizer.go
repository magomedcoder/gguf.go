package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/magomedcoder/gguf.go"
)

// pretokenizePattern - GPT-2 BPE pretokenizer (без lookahead - RE2)
var pretokenizePattern = regexp.MustCompile(`(?i:'s|'t|'re|'ve|'m|'ll|'d)|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+`)

type mergePair struct {
	a, b string
}

// Tokenizer кодирует текст в token IDs и обратно
type Tokenizer struct {
	tokens []string
	id     map[string]int
	merges map[mergePair]int
	bosID  int
	eosID  int
}

// FromGGUF создаёт tokenizer из метаданных GGUF-файла
func FromGGUF(r *gguf.Reader) (*Tokenizer, error) {
	model, err := r.Metadata.String("tokenizer.ggml.model")
	if err != nil {
		return nil, fmt.Errorf("tokenizer: %w", err)
	}

	switch model {
	case "gpt2", "llama":
		return loadBPE(r)
	default:
		return nil, fmt.Errorf("tokenizer: модель %q не поддерживается", model)
	}
}

func loadBPE(r *gguf.Reader) (*Tokenizer, error) {
	tokens, err := r.Metadata.StringArray("tokenizer.ggml.tokens")
	if err != nil {
		return nil, err
	}

	id := make(map[string]int, len(tokens))
	for i, tok := range tokens {
		id[tok] = i
	}

	merges := make(map[mergePair]int)
	if raw, err := r.Metadata.StringArray("tokenizer.ggml.merges"); err == nil {
		for rank, m := range raw {
			parts := strings.SplitN(m, " ", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("tokenizer: некорректный merge %q", m)
			}
			merges[mergePair{parts[0], parts[1]}] = rank
		}
	}

	return &Tokenizer{
		tokens: tokens,
		id:     id,
		merges: merges,
		bosID:  r.Metadata.IntOptional("tokenizer.ggml.bos_token_id", -1),
		eosID:  r.Metadata.IntOptional("tokenizer.ggml.eos_token_id", -1),
	}, nil
}

// BOS возвращает ID токена начала последовательности
func (t *Tokenizer) BOS() int {
	return t.bosID
}

// EOS возвращает ID токена конца последовательности
func (t *Tokenizer) EOS() int {
	return t.eosID
}

// VocabSize возвращает размер словаря
func (t *Tokenizer) VocabSize() int {
	return len(t.tokens)
}

// Encode преобразует текст в token IDs
func (t *Tokenizer) Encode(text string) ([]int, error) {
	if text == "" {
		return nil, nil
	}

	var out []int
	for _, piece := range pretokenizePattern.FindAllString(text, -1) {
		for _, sym := range t.bpe(piece) {
			id, ok := t.id[sym]
			if !ok {
				return nil, fmt.Errorf("tokenizer: неизвестный токен %q", sym)
			}
			out = append(out, id)
		}
	}

	return out, nil
}

// Decode преобразует token IDs в текст
func (t *Tokenizer) Decode(ids []int) string {
	var b strings.Builder
	for _, id := range ids {
		if id < 0 || id >= len(t.tokens) {
			continue
		}
		b.WriteString(t.decodeToken(t.tokens[id]))
	}

	return b.String()
}

func (t *Tokenizer) decodeToken(tok string) string {
	// GPT-2 пробел кодируется префиксом Ġ
	return strings.ReplaceAll(tok, "Ġ", " ")
}

func (t *Tokenizer) bpe(text string) []string {
	if text == "" {
		return nil
	}

	symbols := utf8Symbols(text)
	if len(symbols) == 1 {
		return symbols
	}

	for {
		bestRank := int(^uint(0) >> 1)
		bestIdx := -1

		for i := 0; i < len(symbols)-1; i++ {
			pair := mergePair{symbols[i], symbols[i+1]}
			if rank, ok := t.merges[pair]; ok && rank < bestRank {
				bestRank = rank
				bestIdx = i
			}
		}

		if bestIdx < 0 {
			break
		}

		merged := symbols[bestIdx] + symbols[bestIdx+1]
		next := make([]string, 0, len(symbols)-1)
		next = append(next, symbols[:bestIdx]...)
		next = append(next, merged)
		next = append(next, symbols[bestIdx+2:]...)
		symbols = next
	}

	return symbols
}

func utf8Symbols(s string) []string {
	if s == "" {
		return nil
	}

	n := utf8.RuneCountInString(s)
	out := make([]string, 0, n)
	for _, r := range s {
		out = append(out, string(r))
	}

	return out
}

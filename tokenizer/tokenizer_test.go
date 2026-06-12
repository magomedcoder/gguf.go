package tokenizer

import (
	"testing"

	"github.com/magomedcoder/gguf.go"
)

func TestDecodeGPT2Space(t *testing.T) {
	tok := &Tokenizer{
		tokens: []string{"Ġhello", "Ġworld"},
	}

	if got := tok.Decode([]int{0, 1}); got != " hello world" {
		t.Fatalf("Decode = %q", got)
	}
}

func TestBPEMerge(t *testing.T) {
	tok := &Tokenizer{
		tokens: []string{"he", "llo"},
		id:     map[string]int{"he": 0, "llo": 1, "h": 2, "e": 3, "l": 4, "o": 5},
		merges: map[mergePair]int{
			{"h", "e"}:    0,
			{"he", "l"}:   1,
			{"hel", "l"}:  2,
			{"hell", "o"}: 3,
		},
	}

	symbols := tok.bpe("hello")
	if len(symbols) != 1 || symbols[0] != "hello" {
		t.Fatalf("bpe(hello) = %v", symbols)
	}
}

func TestFromGGUFMetadata(t *testing.T) {
	r := &gguf.Reader{
		Metadata: gguf.Metadata{
			"tokenizer.ggml.model":        "gpt2",
			"tokenizer.ggml.tokens":       []string{"<s>", "</s>", "Ġhi"},
			"tokenizer.ggml.merges":       []string{"Ġ h", "h i"},
			"tokenizer.ggml.bos_token_id": int32(0),
			"tokenizer.ggml.eos_token_id": int32(1),
		},
	}

	tok, err := FromGGUF(r)
	if err != nil {
		t.Fatalf("FromGGUF: %v", err)
	}

	if tok.BOS() != 0 || tok.EOS() != 1 {
		t.Fatalf("special tokens: bos=%d eos=%d", tok.BOS(), tok.EOS())
	}

	if tok.VocabSize() != 3 {
		t.Fatalf("vocab size = %d", tok.VocabSize())
	}
}

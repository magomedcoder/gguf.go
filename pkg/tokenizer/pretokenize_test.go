package tokenizer

import (
	"testing"
)

func TestByteRoundtrip(t *testing.T) {
	for _, s := range []string{"Hello", "\n", " "} {
		enc := byteEncodeWord(s)
		dec := byteDecodeToken(enc)
		if dec != s {
			t.Fatalf("%q: enc=%q dec=%q", s, enc, dec)
		}
	}
}

func TestPretokenizeQwen2Newline(t *testing.T) {
	parts := pretokenizeQwen2("a\nb")
	if len(parts) < 2 {
		t.Fatalf("parts = %v", parts)
	}
}

func TestEncodeNewlineToken(t *testing.T) {
	tok := testTokenizer(t)

	ids, err := tok.Encode("\n")
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 1 {
		t.Fatalf("newline ids = %v", ids)
	}
}

func testTokenizer(t *testing.T) *Tokenizer {
	t.Helper()

	merges := map[mergePair]int{}
	// минимальный словарный запас с переводом строки в виде Ċ (байт 10)
	bs := byteToUnicode[10]
	tokens := []string{"<s>", "</s>", bs, "a", "b"}
	id := map[string]int{}
	for i, s := range tokens {
		id[s] = i
	}

	return &Tokenizer{
		tokens:      tokens,
		id:          id,
		merges:      merges,
		byteEncode:  true,
		pretokenize: pretokenizeQwen2,
	}
}

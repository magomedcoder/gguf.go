package tokenizer

import (
	"strings"
	"unicode/utf8"
)

var (
	byteToUnicode map[byte]string
	unicodeToByte map[string]byte
)

func init() {
	byteToUnicode = make(map[byte]string, 256)
	unicodeToByte = make(map[string]byte, 256)

	add := func(b byte, s string) {
		byteToUnicode[b] = s
		unicodeToByte[s] = b
	}

	for ch := 0x21; ch <= 0x7E; ch++ {
		add(byte(ch), string(rune(ch)))
	}

	for ch := 0xA1; ch <= 0xAC; ch++ {
		add(byte(ch), string(rune(ch)))
	}

	for ch := 0xAE; ch <= 0xFF; ch++ {
		add(byte(ch), string(rune(ch)))
	}

	n := 0
	for ch := range 256 {
		if _, ok := byteToUnicode[byte(ch)]; ok {
			continue
		}

		add(byte(ch), string(rune(256+n)))
		n++
	}
}

func byteEncodeWord(word string) string {
	var b strings.Builder
	for i := 0; i < len(word); i++ {
		b.WriteString(byteToUnicode[word[i]])
	}

	return b.String()
}

func byteDecodeToken(tok string) string {
	out := make([]byte, 0, len(tok))
	for i := 0; i < len(tok); {
		r, size := utf8.DecodeRuneInString(tok[i:])
		if b, ok := unicodeToByte[string(r)]; ok {
			out = append(out, b)
		}

		i += size
	}

	return string(out)
}

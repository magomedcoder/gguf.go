package runtime

import (
	"github.com/magomedcoder/gguf.go"
	"github.com/magomedcoder/gguf.go/model"
	"github.com/magomedcoder/gguf.go/tokenizer"
)

// Engine загружает GGUF-модель для inference
type Engine struct {
	Model model.Model
	tok   *tokenizer.Tokenizer
}

// Load открывает GGUF-файл и загружает модель
func Load(path string) (*Engine, error) {
	r, err := gguf.OpenFile(path)
	if err != nil {
		return nil, err
	}

	m, err := model.Load(r)
	if err != nil {
		return nil, err
	}

	tok, err := tokenizer.FromGGUF(r)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Model: m,
		tok:   tok,
	}, nil
}

// Tokenizer возвращает tokenizer модели
func (e *Engine) Tokenizer() *tokenizer.Tokenizer {
	return e.tok
}

// ForwardTokenIDs выполняет forward pass для token IDs
func (e *Engine) ForwardTokenIDs(tokens []int, startPos int) ([]float32, error) {
	return e.Model.Forward(tokens, startPos)
}

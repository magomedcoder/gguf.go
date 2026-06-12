package runtime

import (
	"fmt"
	"io"

	"github.com/magomedcoder/gguf.go/sampler"
	"github.com/magomedcoder/gguf.go/tokenizer"
)

// GenerateParams - параметры генерации
type GenerateParams struct {
	MaxTokens int
	Sampler   sampler.Func
	OnToken   func(tokenID int) bool
}

// Context выполняет prefill и autoregressive decode
type Context struct {
	engine *Engine
	tok    *tokenizer.Tokenizer
}

// NewContext создаёт inference-контекст
func (e *Engine) NewContext() (*Context, error) {
	if e.tok == nil {
		return nil, fmt.Errorf("runtime: tokenizer не загружен")
	}

	return &Context{
		engine: e,
		tok:    e.tok,
	}, nil
}

// Encode преобразует текст в token IDs
func (c *Context) Encode(text string) ([]int, error) {
	return c.tok.Encode(text)
}

// Generate выполняет prefill + decode и возвращает сгенерированный текст
func (c *Context) Generate(prompt string, params GenerateParams) (string, error) {
	if params.Sampler == nil {
		params.Sampler = sampler.Greedy
	}

	if params.MaxTokens <= 0 {
		params.MaxTokens = 128
	}

	c.engine.Model.ResetCache()

	promptTokens, err := c.tok.Encode(prompt)
	if err != nil {
		return "", err
	}

	logits, err := c.engine.Model.Forward(promptTokens, 0)
	if err != nil {
		return "", err
	}

	var generated []int
	eos := c.tok.EOS()

	for i := 0; i < params.MaxTokens; i++ {
		next := params.Sampler(logits)
		if next < 0 {
			break
		}
		
		if eos >= 0 && next == eos {
			break
		}

		generated = append(generated, next)
		if params.OnToken != nil && !params.OnToken(next) {
			break
		}

		startPos := len(promptTokens) + len(generated) - 1
		logits, err = c.engine.Model.Forward([]int{next}, startPos)
		if err != nil {
			return "", err
		}
	}

	return c.tok.Decode(generated), nil
}

// GenerateStream пишет сгенерированные token IDs в w по мере decode
func (c *Context) GenerateStream(prompt string, params GenerateParams, w io.Writer) error {
	params.OnToken = func(id int) bool {
		_, err := io.WriteString(w, c.tok.Decode([]int{id}))
		return err == nil
	}

	_, err := c.Generate(prompt, params)
	return err
}

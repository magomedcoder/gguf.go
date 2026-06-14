package gguf

import (
	"io"

	"github.com/magomedcoder/gguf.go/pkg/format"
	"github.com/magomedcoder/gguf.go/pkg/runtime"
	"github.com/magomedcoder/gguf.go/pkg/sampler"
)

// Парсинг GGUF (pkg/format)

type (
	Reader     = format.Reader
	Metadata   = format.Metadata
	TensorInfo = format.TensorInfo
	Type       = format.Type
	GGML       = format.GGML
	Filetype   = format.Filetype
)

// OpenFile открывает GGUF-файл по пути на диске
func OpenFile(filename string) (*Reader, error) {
	return format.OpenFile(filename)
}

// Open парсит GGUF из потока для чтения тензоров источник должен реализовать io.ReaderAt
func Open(readSeeker io.ReadSeeker) (*Reader, error) {
	return format.Open(readSeeker)
}

// Inference (pkg/runtime, pkg/sampler)

type (
	Engine         = runtime.Engine
	Context        = runtime.Context
	GenerateParams = runtime.GenerateParams
	SamplerFunc    = sampler.Func
	SamplerConfig  = sampler.Config
)

// Load загружает модель и tokenizer из GGUF-файла
func Load(path string) (*Engine, error) {
	return runtime.Load(path)
}

// NewSampler возвращает функцию выбора следующего токена (greedy, temperature, top-k, top-p)
func NewSampler(cfg SamplerConfig) SamplerFunc {
	return sampler.New(cfg)
}

// Greedy выбирает токен с максимальным logit
var Greedy = sampler.Greedy

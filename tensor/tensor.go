package tensor

import (
	"fmt"
	"io"

	"github.com/magomedcoder/gguf.go"
	"github.com/magomedcoder/gguf.go/quant"
)

// Tensor - деквантизованный тензор в float32
type Tensor struct {
	Name  string
	Shape []int
	Type  gguf.GGML
	Data  []float32
	Raw   []byte
}

// LoadRaw читает сырые байты тензора из GGUF-файла
func LoadRaw(info *gguf.TensorInfo) ([]byte, error) {
	r, err := info.Reader()
	if err != nil {
		return nil, err
	}

	raw := make([]byte, info.Size())
	if _, err := io.ReadFull(r, raw); err != nil {
		return nil, fmt.Errorf("tensor %q: %w", info.Name, err)
	}

	return raw, nil
}

// LoadFloats загружает и деквантизирует тензор в float32
func LoadFloats(info *gguf.TensorInfo) ([]float32, error) {
	raw, err := LoadRaw(info)
	if err != nil {
		return nil, err
	}

	return quant.ToFloat32(info.Type, raw, int(info.ValuesCount()))
}

// FromGGUF загружает тензор из GGUF
func FromGGUF(info *gguf.TensorInfo) (*Tensor, error) {
	raw, err := LoadRaw(info)
	if err != nil {
		return nil, err
	}

	n := int(info.ValuesCount())
	data, err := quant.ToFloat32(info.Type, raw, n)
	if err != nil {
		return nil, err
	}

	shape := make([]int, len(info.Dimensions))
	for i, d := range info.Dimensions {
		shape[i] = int(d)
	}

	return &Tensor{
		Name:  info.Name,
		Shape: shape,
		Type:  info.Type,
		Data:  data,
		Raw:   raw,
	}, nil
}

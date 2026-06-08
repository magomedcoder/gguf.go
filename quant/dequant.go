package quant

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/magomedcoder/gguf.go"
)

// ToFloat32 деквантизирует сырые байты GGML-тензора в float32
func ToFloat32(typ gguf.GGML, data []byte, n int) ([]float32, error) {
	switch typ {
	case gguf.GgmlFloat32:
		return dequantF32(data, n)
	case gguf.GgmlFloat16:
		return dequantF16(data, n)
	case gguf.GgmlQ8_0:
		return DequantQ8_0(data, n)
	case gguf.GgmlInt32:
		return dequantI32(data, n)
	default:
		return nil, fmt.Errorf("quant: тип %s пока не поддерживается", typ)
	}
}

// dequantF32 читает n float32 из буфера
func dequantF32(data []byte, n int) ([]float32, error) {
	want := n * 4
	if len(data) < want {
		return nil, fmt.Errorf("quant: F32 данных недостаточно: нужно %d, есть %d", want, len(data))
	}

	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:]))
	}

	return out, nil
}

// dequantF16 конвертирует n fp16 в float32
func dequantF16(data []byte, n int) ([]float32, error) {
	want := n * 2
	if len(data) < want {
		return nil, fmt.Errorf("quant: F16 данных недостаточно: нужно %d, есть %d", want, len(data))
	}

	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = FP16ToFP32(binary.LittleEndian.Uint16(data[i*2:]))
	}

	return out, nil
}

// dequantI32 читает n int32 и приводит к float32
func dequantI32(data []byte, n int) ([]float32, error) {
	want := n * 4
	if len(data) < want {
		return nil, fmt.Errorf("quant: I32 данных недостаточно: нужно %d, есть %d", want, len(data))
	}

	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = float32(int32(binary.LittleEndian.Uint32(data[i*4:])))
	}

	return out, nil
}

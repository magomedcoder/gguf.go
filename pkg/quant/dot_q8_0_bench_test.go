package quant

import "testing"

func BenchmarkDotBlockQ8_0(b *testing.B) {
	var block [BlockQ8_0Size]byte
	block[2] = 1
	block[3] = 2
	x := make([]float32, QK8_0)
	for i := range x {
		x[i] = 0.01
	}

	b.ResetTimer()
	for range b.N {
		_, _ = DotBlockQ8_0(block[:], x)
	}
}

//go:build amd64

package quant

import "encoding/binary"

func init() {
	if cpuHasAVX2() {
		dotBlockQ8_0 = dotBlockQ8_0AVX2
	}
}

func dotBlockQ8_0AVX2(block []byte, x []float32) float32 {
	d := FP16ToFP32(binary.LittleEndian.Uint16(block[0:2]))
	return dotBlockQ8_0AVX2Asm(d, &block[2], &x[0])
}

//go:noescape
func cpuHasAVX2() bool

//go:noescape
func dotBlockQ8_0AVX2Asm(d float32, q *byte, x *float32) float32

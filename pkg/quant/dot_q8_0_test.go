package quant

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestDotBlockQ8_0MatchesGo(t *testing.T) {
	var block [BlockQ8_0Size]byte
	binary.LittleEndian.PutUint16(block[0:2], 0x3800) // scale 0.5
	for i := range QK8_0 {
		block[2+i] = byte(int8((i % 7) - 3))
	}

	x := make([]float32, QK8_0)
	for i := range x {
		x[i] = float32(i)*0.1 - 0.5
	}

	gotFn, err := DotBlockQ8_0(block[:], x)
	if err != nil {
		t.Fatal(err)
	}

	want := dotBlockQ8_0Go(block[:], x)
	if math.Abs(float64(gotFn-want)) > 1e-2 {
		t.Fatalf("dot = %v, dotGo = %v", gotFn, want)
	}
}

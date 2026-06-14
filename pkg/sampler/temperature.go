package sampler

import (
	"math"
	"math/rand/v2"
)

// Func выбирает следующий token ID из logits
type Func func(logits []float32) int

// GreedyFunc возвращает greedy sampler
func GreedyFunc() Func {
	return Greedy
}

// Temperature возвращает sampler с температурой; temp <= 0 эквивалентен greedy
func Temperature(temp float32, rng *rand.Rand) Func {
	if temp <= 0 {
		return Greedy
	}

	return func(logits []float32) int {
		if len(logits) == 0 {
			return -1
		}

		maxLogit := logits[0]
		for _, v := range logits[1:] {
			if v > maxLogit {
				maxLogit = v
			}
		}

		probs := make([]float64, len(logits))
		sum := 0.0
		invTemp := 1.0 / float64(temp)
		for i, v := range logits {
			p := math.Exp(float64(v-maxLogit) * invTemp)
			probs[i] = p
			sum += p
		}

		r := rng.Float64() * sum
		for i, p := range probs {
			r -= p
			if r <= 0 {
				return i
			}
		}

		return len(logits) - 1
	}
}

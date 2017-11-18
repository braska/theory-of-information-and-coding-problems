package entropy

import "math"

func Entropy(countSymbols int, symbols map[string]int) (entropy float64) {
	entropy = 0.0
	for _, value := range symbols {
		p_i := float64(value) / float64(countSymbols)
		entropy -= p_i * math.Log2(p_i)
	}
	return
}

func ConditionalEntropy(countSymbols int, symbols map[string]int, combinations map[string]int, countCombinations int) (conditionalEntropy float64) {
	conditionalEntropy = 0.0
	for s1, count1 := range symbols {
		p_i := float64(count1) / float64(countSymbols)
		for s2 := range symbols {
			countForCombination, ok := combinations[s1+s2]
			if ok {
				p_i_j := float64(countForCombination) / float64(countCombinations)
				conditionalEntropy -= p_i_j * math.Log2(p_i_j/p_i)
			}
		}
	}
	return
}
package entity

import "math"

type Score int

func NewScore(debtAmounts []float64) Score {
	var sum float64
	for _, v := range debtAmounts {
		sum += v
	}
	if length := len(debtAmounts); length > 0 {
		avg := sum / float64(length)
		return calculateScore(avg)
	}
	return calculateScore(sum)
}

func calculateScore(x float64) Score {
	return Score(10000 / math.Sqrt(x+100))
}

package entity

import (
	"testing"
)

func TestNewScore(t *testing.T) {
	t.Run("teste", func(t *testing.T) {
		debtAmounts := []float64{1000.0}
		result := NewScore(debtAmounts)
		t.Error(result)
	})
}

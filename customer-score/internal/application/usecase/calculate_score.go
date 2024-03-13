package usecase

import (
	"context"
	"github.com/carloseduribeiro/auth-challenge/customer-score/internal/domain/entity"
)

type CalculateScore struct {
	repository entity.DebtsRepository
}

type ScoreOutputDto struct {
	Score int `json:"score"`
}

func (c *CalculateScore) Execute(ctx context.Context, document string) (*ScoreOutputDto, error) {
	debts, err := c.repository.GetDebtsByDocument(ctx, document)
	if err != nil {
		return nil, err
	}
	debtAmounts := make([]float64, 0, len(debts))
	for _, d := range debts {
		debtAmounts = append(debtAmounts, d.Amount)
	}
	return &ScoreOutputDto{Score: int(entity.NewScore(debtAmounts))}, nil
}

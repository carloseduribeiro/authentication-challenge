package usecase

import (
	"context"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

type GetDebts struct {
	repository entity.DebtsRepository
}

func NewGetDebtsUseCase(repository entity.DebtsRepository) *GetDebts {
	return &GetDebts{
		repository: repository,
	}
}

type GetDebtsOutputDto struct {
	Id        uuid.UUID
	Document  string
	Amount    float64
	DueDate   time.Time
	CreatedAt time.Time
}

func (g *GetDebts) Execute(ctx context.Context, document string) ([]GetDebtsOutputDto, error) {
	debts, err := g.repository.FindAll(ctx, document)
	if err != nil {
		return nil, err
	}
	result := make([]GetDebtsOutputDto, 0, len(debts))
	for _, debt := range debts {
		result = append(result, GetDebtsOutputDto{
			Id:        debt.Id(),
			Document:  debt.Document(),
			Amount:    debt.Amount(),
			DueDate:   debt.DueDate(),
			CreatedAt: debt.CreatedAt(),
		})
	}
	return result, nil
}

package usecase

import (
	"context"
	"github.com/carloseduribeiro/authentication-challenge/customer-debts/internal/domain/entity"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/date"
	"github.com/google/uuid"
	"time"
)

type CreateDebt struct {
	repository entity.DebtsRepository
}

func NewCreateDebtUseCase(repository entity.DebtsRepository) *CreateDebt {
	return &CreateDebt{
		repository: repository,
	}
}

type CreateDebtInputDto struct {
	Document string    `json:"cpf"`
	Amount   float64   `json:"valor"`
	DueDate  date.Date `json:"vencimento"`
}

type CreateDebtOutputDto struct {
	Id       uuid.UUID `json:"id"`
	Document string    `json:"cpf"`
	Amount   float64   `json:"valor"`
	DueDate  date.Date `json:"vencimento"`
}

func (c *CreateDebt) Execute(ctx context.Context, input *CreateDebtInputDto) (*CreateDebtOutputDto, error) {
	id, _ := uuid.NewUUID()
	debt, err := entity.NewDebt(id, input.Document, input.Amount, input.DueDate.T, time.Now())
	if err != nil {
		return nil, err
	}
	if err = c.repository.Save(ctx, debt); err != nil {
		return nil, err
	}
	return &CreateDebtOutputDto{
		Id:       id,
		Document: input.Document,
		Amount:   input.Amount,
		DueDate:  input.DueDate,
	}, nil
}

package database

import (
	"context"
	"github.com/carloseduribeiro/authentication-challenge/customer-debts/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DebtsRepository struct {
	q *Queries
}

func NewDebtsRepository(dbPool *pgxpool.Pool) *DebtsRepository {
	return &DebtsRepository{q: New(dbPool)}
}

func (d *DebtsRepository) FindAll(ctx context.Context, document string) ([]entity.Debt, error) {
	debts, err := d.q.GetDebtsByDocument(ctx, document)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Debt, 0, len(debts))
	for _, debt := range debts {
		e, _ := entity.NewDebt(debt.ID, debt.Document, debt.Amount, debt.Duedate, debt.CreatedAt)
		result = append(result, *e)
	}
	return result, nil
}

func (d *DebtsRepository) Save(ctx context.Context, debt *entity.Debt) error {
	params := InsertDebtParams{
		ID:        debt.Id(),
		Document:  debt.Document(),
		Duedate:   debt.DueDate(),
		Amount:    debt.Amount(),
		CreatedAt: debt.CreatedAt(),
	}
	return d.q.InsertDebt(ctx, params)
}

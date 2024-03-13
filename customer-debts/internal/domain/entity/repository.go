package entity

import "context"

type DebtsRepository interface {
	FindAll(ctx context.Context, document string) ([]Debt, error)
	Save(ctx context.Context, debt *Debt) error
}

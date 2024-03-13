package entity

import "context"

type DebtsRepository interface {
	GetDebtsByDocument(ctx context.Context, document string) ([]Debt, error)
}

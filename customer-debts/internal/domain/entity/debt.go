package entity

import (
	"errors"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/cpf"
	"github.com/google/uuid"
	"time"
)

type Debt struct {
	id        uuid.UUID
	document  string
	amount    float64
	dueDate   time.Time
	createdAt time.Time
}

func NewDebt(id uuid.UUID, document string, amount float64, dueDate, createdAt time.Time) (*Debt, error) {
	if !cpf.Validate(document) {
		return nil, errors.New("invalid cpf")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	return &Debt{
		id:        id,
		document:  document,
		amount:    amount,
		dueDate:   dueDate,
		createdAt: createdAt,
	}, nil
}

func (d Debt) Id() uuid.UUID {
	return d.id
}

func (d Debt) Document() string {
	return d.document
}
func (d Debt) Amount() float64 {
	return d.amount
}
func (d Debt) DueDate() time.Time {
	return d.dueDate
}
func (d Debt) CreatedAt() time.Time {
	return d.createdAt
}

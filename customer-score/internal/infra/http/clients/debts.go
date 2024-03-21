package clients

import (
	"context"
	"encoding/json"
	"github.com/carloseduribeiro/authentication-challenge/customer-score/internal/domain/entity"
	"net/http"
)

const getDebtsResource = "/customer/debts/"

type Debts struct {
	serviceHost string
}

func (d *Debts) GetDebtsByDocument(ctx context.Context, document string) ([]entity.Debt, error) {
	response, err := http.Get(d.serviceHost + getDebtsResource + document)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var result []entity.Debt
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

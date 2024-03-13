package handlers

import (
	"encoding/json"
	"errors"
	"github.com/carloseduribeiro/authentication-challenge/customer-debts/internal/application/usecase"
	"github.com/carloseduribeiro/authentication-challenge/customer-debts/internal/domain/entity"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/cpf"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type CreateDebt struct {
	repository entity.DebtsRepository
}

func NewCreateDebtHandler(repository entity.DebtsRepository) *CreateDebt {
	return &CreateDebt{repository: repository}
}

func (c *CreateDebt) Handle(w http.ResponseWriter, r *http.Request) {
	requestBody := new(usecase.CreateDebtInputDto)
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: "error parsing request body",
			Errors:  []string{err.Error()},
		})
		return
	}
	if errorS := validateCreateDebtPayload(requestBody); len(errorS) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: "invalid parameters on request body",
			Errors:  errorS,
		})
		return
	}
	useCase := usecase.NewCreateDebtUseCase(c.repository)
	result, err := useCase.Execute(r.Context(), requestBody)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Message: "error creating debt",
				Errors:  []string{err.Error()},
			})
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
	return
}

// validateCreateDebtPayload validate the input fields and returns a slice with validation error messages
func validateCreateDebtPayload(input *usecase.CreateDebtInputDto) []string {
	errS := make([]string, 0, 5)
	if !cpf.Validate(input.Document) {
		errS = append(errS, "invalid cpf")
	}
	if input.Amount <= 0 {
		errS = append(errS, "amount must be greater than zero")
	}
	return errS
}

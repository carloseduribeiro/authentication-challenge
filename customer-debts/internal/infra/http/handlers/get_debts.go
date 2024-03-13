package handlers

import (
	"encoding/json"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/internal/application/usecase"
	"github.com/carloseduribeiro/auth-challenge/customer-debts/internal/domain/entity"
	"github.com/carloseduribeiro/auth-challenge/lib-utils/pkg/cpf"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type GetDebts struct {
	repository entity.DebtsRepository
}

func NewGetDebtsHandler(repository entity.DebtsRepository) *GetDebts {
	return &GetDebts{
		repository: repository,
	}
}

func (h *GetDebts) Handle(w http.ResponseWriter, r *http.Request) {
	document := chi.URLParam(r, "userDocument")
	if !cpf.Validate(document) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	debts, err := usecase.NewGetDebtsUseCase(h.repository).Execute(r.Context(), document)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(debts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(debts)
}

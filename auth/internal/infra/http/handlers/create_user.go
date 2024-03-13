package handlers

import (
	"encoding/json"
	"errors"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/application/usecase"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/domain/entity"
	"github.com/google/uuid"
	"net/http"
)

type CreateUser struct {
	repository        entity.UserRepository
	uuidGeneratorFunc func() (uuid.UUID, error)
}

func NewCreateUser(repository entity.UserRepository, uuidGeneratorFunc func() (uuid.UUID, error)) *CreateUser {
	return &CreateUser{
		repository:        repository,
		uuidGeneratorFunc: uuidGeneratorFunc,
	}
}

func (c *CreateUser) Handler(w http.ResponseWriter, r *http.Request) {
	requestBody := new(usecase.CreateUserInputDto)
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errorS := validateCreateUserPayload(requestBody); len(errorS) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: "invalid parameters on request body",
			Errors:  errorS,
		})
		return
	}
	useCase := usecase.NewCreateUserUseCase(c.repository, c.uuidGeneratorFunc)
	result, err := useCase.Execute(r.Context(), *requestBody)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{
				Message: "error creating user",
				Errors:  []string{err.Error()},
			})
			return
		} else {
			// TODO: implement log with traces
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
	return
}

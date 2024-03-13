package handlers

import (
	"encoding/json"
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/application/usecase"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"net/http"
	"time"
)

type Login struct {
	userRepository    entity.UserRepository
	sessionRepository entity.SessionRepository
	tokenSecretKey    string
	sessionDuration   time.Duration
}

func NewLogin(userRepository entity.UserRepository, sessionRepository entity.SessionRepository, tokenSecretKey string, sessionDuration time.Duration) *Login {
	return &Login{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		tokenSecretKey:    tokenSecretKey,
		sessionDuration:   sessionDuration,
	}
}

func (l *Login) Handler(w http.ResponseWriter, r *http.Request) {
	requestBody := new(usecase.LoginInputDto)
	if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	useCase := usecase.NewLoginUseCase(l.userRepository, l.sessionRepository, l.tokenSecretKey, l.sessionDuration)
	result, err := useCase.Execute(r.Context(), requestBody)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(result)
	return
}

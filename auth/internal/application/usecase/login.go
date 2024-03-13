package usecase

import (
	"context"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/domain/entity"
	"github.com/google/uuid"
	"time"
)

const DefaultTokenType = "Bearer"

type SessionConfig struct {
	sessionDuration time.Duration
	secretKey       string
}

type Login struct {
	userRepository    entity.UserRepository
	sessionRepository entity.SessionRepository
	sessionConfig     *SessionConfig
}

func NewLoginUseCase(repository entity.UserRepository, sessionRepository entity.SessionRepository, tokenSecretKey string, sessionDuration time.Duration) *Login {
	return &Login{
		userRepository:    repository,
		sessionRepository: sessionRepository,
		sessionConfig: &SessionConfig{
			secretKey:       tokenSecretKey,
			sessionDuration: sessionDuration,
		},
	}
}

type LoginInputDto struct {
	Document string `json:"cpf"`
	Password string `json:"senha"`
}

type LoginOutputDto struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Type    string `json:"type"`
}

func (l *Login) Execute(ctx context.Context, input *LoginInputDto) (*LoginOutputDto, error) {
	u, err := l.userRepository.GetUserByDocument(ctx, input.Document)
	if err != nil {
		return nil, err
	}
	sessionId, _ := uuid.NewUUID()
	token, err := u.NewSession(sessionId, time.Now(), l.sessionConfig.sessionDuration, l.sessionConfig.secretKey)
	if err != nil {
		return nil, err
	}
	if err = l.sessionRepository.Save(ctx, u.ActiveSession()); err != nil {
		return nil, err
	}
	return &LoginOutputDto{
		Message: "login successful",
		Token:   token,
		Type:    DefaultTokenType,
	}, nil
}

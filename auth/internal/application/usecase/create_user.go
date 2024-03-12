package usecase

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity/user"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"github.com/google/uuid"
	"time"
)

type uuidGeneratorFunc func() (uuid.UUID, error)

type CreateUser struct {
	repository        user.Repository
	uuidGeneratorFunc uuidGeneratorFunc
}

func NewCreateUserUseCase(repository user.Repository, uuidGeneratorFunc uuidGeneratorFunc) *CreateUser {
	return &CreateUser{
		repository:        repository,
		uuidGeneratorFunc: uuidGeneratorFunc,
	}
}

type CreateUserInputDto struct {
	Document  string    `json:"cpf"`
	Name      string    `json:"nome"`
	BirthDate time.Time `json:"nascimento"`
	Email     string    `json:"email"`
	Password  string    `json:"senha"`
}

type CreatedUserOutputDto struct {
	ID uuid.UUID `json:"id"`
	*CreateUserInputDto
}

func isErrUserNotFound(err error) bool {
	return errors.Is(err, database.ErrUserNotFound)
}

func (c *CreateUser) Execute(ctx context.Context, input CreateUserInputDto) (*CreatedUserOutputDto, error) {
	if u, err := c.repository.GetUserByDocument(ctx, input.Document); err != nil && !isErrUserNotFound(err) {
		return nil, err
	} else if u != nil {
		return nil, errors.New("user already exists with the given document")
	}
	if u, err := c.repository.GetUserByEmail(ctx, input.Email); err != nil && !isErrUserNotFound(err) {
		return nil, err
	} else if u != nil {
		return nil, errors.New("user already exists with the given email")
	}
	u, err := user.NewUser(
		input.Document, input.Name, input.Email, input.Password, input.BirthDate, user.WithUUIDGeneratorFunc(c.uuidGeneratorFunc))
	if err != nil {
		return nil, err
	}
	if err = c.repository.Create(ctx, u); err != nil {
		return nil, err
	}
	input.Password = u.Password()
	return &CreatedUserOutputDto{
		ID:                 u.ID(),
		CreateUserInputDto: &input,
	}, nil
}

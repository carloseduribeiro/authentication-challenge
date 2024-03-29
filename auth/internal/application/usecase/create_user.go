package usecase

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/domain/entity"
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/infra/database"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/date"
	"github.com/google/uuid"
)

type CreateUser struct {
	repository        entity.UserRepository
	uuidGeneratorFunc func() (uuid.UUID, error)
}

func NewCreateUserUseCase(repository entity.UserRepository, uuidGeneratorFunc func() (uuid.UUID, error)) *CreateUser {
	return &CreateUser{
		repository:        repository,
		uuidGeneratorFunc: uuidGeneratorFunc,
	}
}

type CreateUserInputDto struct {
	Document  string    `json:"cpf"`
	Name      string    `json:"nome"`
	BirthDate date.Date `json:"nascimento"`
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

var ErrUserAlreadyExists = errors.New("user already exists")

func (c *CreateUser) Execute(ctx context.Context, input CreateUserInputDto) (*CreatedUserOutputDto, error) {
	if u, err := c.repository.GetUserByDocument(ctx, input.Document); err != nil && !isErrUserNotFound(err) {
		return nil, err
	} else if u != nil {
		return nil, ErrUserAlreadyExists
	}
	if u, err := c.repository.GetUserByEmail(ctx, input.Email); err != nil && !isErrUserNotFound(err) {
		return nil, err
	} else if u != nil {
		return nil, ErrUserAlreadyExists
	}
	id, err := c.uuidGeneratorFunc()
	if err != nil {
		return nil, err
	}
	u, err := entity.NewUser(
		id, input.Document, input.Name, input.Email, input.BirthDate.T, entity.WithPassword(input.Password),
	)
	if err != nil {
		return nil, err
	}
	if err = c.repository.Create(ctx, u); err != nil {
		return nil, err
	}
	input.Password = u.Password()
	return &CreatedUserOutputDto{
		ID:                 id,
		CreateUserInputDto: &input,
	}, nil
}

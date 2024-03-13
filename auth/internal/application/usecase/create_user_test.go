package usecase

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	entityMocks "github.com/carloseduribeiro/auth-challenge/auth/mocks/internal_/domain/entity"
	"github.com/carloseduribeiro/auth-challenge/auth/pkg/date"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var (
	staticUUID    = uuid.New()
	validDocument = "72430024985"
)

func fakeUUIDGenerator() (uuid.UUID, error) {
	return staticUUID, nil
}

type CreateUserTestSuite struct {
	repoMock *entityMocks.UserRepository
	useCase  *CreateUser
	suite.Suite
}

func (c *CreateUserTestSuite) SetupSubTest() {
	c.repoMock = entityMocks.NewUserRepository(c.T())
	c.useCase = NewCreateUserUseCase(c.repoMock, fakeUUIDGenerator)
}

func (c *CreateUserTestSuite) TestExecute() {
	c.Run("must return an error from repository when GetUserByDocument method call returns error", func() {
		// given
		input := CreateUserInputDto{}
		someErr := errors.New("some error")
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, someErr)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.Error(err)
		c.ErrorIs(err, someErr)
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("must return an user already exists with document error", func() {
		// given
		input := CreateUserInputDto{}
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(&entity.User{}, nil)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.ErrorIs(err, ErrUserAlreadyExists)
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("must return an error from repository when GetUserByEmail method call returns error", func() {
		// given
		input := CreateUserInputDto{}
		someErr := errors.New("some error")
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, someErr)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.Error(err)
		c.ErrorIs(err, someErr)
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("must return an user already exists with document email", func() {
		// given
		input := CreateUserInputDto{}
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(&entity.User{}, nil)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.ErrorIs(err, ErrUserAlreadyExists)
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("must return an error from entity when NewUser method call returns error", func() {
		// given
		invalidDocument := "12345"
		input := CreateUserInputDto{Document: invalidDocument}
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.Error(err)
		c.EqualError(err, "invalid cpf")
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("must return an error from repository when Create method call returns error", func() {
		// given
		input := CreateUserInputDto{
			Document:  validDocument,
			Name:      "Jhon",
			BirthDate: date.New(time.Now()),
			Email:     "jhon@doe.com",
			Password:  "password",
		}
		fakeErr := pgx.ErrTxCommitRollback
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(fakeErr)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.Nil(result)
		c.Error(err)
		c.ErrorIs(err, fakeErr)
		c.repoMock.AssertExpectations(c.T())
	})

	c.Run("should return the new user created", func() {
		// given
		input := CreateUserInputDto{
			Document:  validDocument,
			Name:      "Jhon",
			BirthDate: date.New(time.Now()),
			Email:     "jhon@doe.com",
			Password:  "password",
		}
		// when
		c.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		c.repoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
		result, err := c.useCase.Execute(context.TODO(), input)
		// then
		c.NoError(err)
		c.Equal(staticUUID.String(), result.ID.String())
		c.Equal(input.Document, result.Document)
		c.Equal(input.Name, result.Name)
		c.Equal(input.BirthDate, result.BirthDate)
		c.Equal(input.Email, result.Email)
		c.NoError(bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(input.Password)))
		c.repoMock.AssertExpectations(c.T())
	})
}

func TestCreateUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserTestSuite))
}

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/application/usecase"
	entity "github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity/user"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/infra/database"
	"github.com/carloseduribeiro/auth-challenge/auth/mocks/internal_/domain/entity/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CreateUserHandlerTestSuite struct {
	repoMock          *user.Repository
	createUserHandler *CreateUser
	rr                *httptest.ResponseRecorder
	fakeId            uuid.UUID
	uuidGeneratorFunc func() (uuid.UUID, error)
	suite.Suite
}

func (t *CreateUserHandlerTestSuite) SetupTest() {
	t.fakeId = uuid.New()
	t.uuidGeneratorFunc = func() (uuid.UUID, error) {
		return t.fakeId, nil
	}
	t.createUserHandler = &CreateUser{uuidGeneratorFunc: t.uuidGeneratorFunc}
}

func (t *CreateUserHandlerTestSuite) SetupSubTest() {
	t.repoMock = user.NewRepository(t.T())
	t.rr = httptest.NewRecorder()
	t.createUserHandler.repository = t.repoMock
}

func (t *CreateUserHandlerTestSuite) TestHandler() {
	t.Run("must return http status BadRequest with no response body when got an invalid request body", func() {
		// given
		req, err := http.NewRequest(http.MethodPost, "/auth/users", strings.NewReader(``))
		t.Require().NoError(err)
		// when
		handlerFunc := http.HandlerFunc(t.createUserHandler.Handler)
		handlerFunc.ServeHTTP(t.rr, req)
		// then
		t.Equal(http.StatusBadRequest, t.rr.Code)
		t.Empty(t.rr.Body)
		t.repoMock.AssertExpectations(t.T())
	})

	t.Run("must return http status BadRequest with invalid request body error response", func() {
		// given
		req, err := http.NewRequest(http.MethodPost, "/auth/users", strings.NewReader(`{}`))
		t.Require().NoError(err)
		expectedResponseBody := `{
			"message": "invalid parameters on request body",
			"errors": [
				"invalid cpf",
				"name must be at least 3 characters",
				"invalid email address",
				"password must contains at least one lower case character",
				"password must contains at least one upper case character",
				"password must contains at least one numeric character",
				"password must contains at least 6 or up to 12 characters"
			]
		}`
		// when
		handlerFunc := http.HandlerFunc(t.createUserHandler.Handler)
		handlerFunc.ServeHTTP(t.rr, req)
		// then
		t.Equal(http.StatusBadRequest, t.rr.Code)
		t.JSONEq(expectedResponseBody, t.rr.Body.String())
		t.repoMock.AssertExpectations(t.T())
	})

	t.Run("must return http status BadRequest and the error when usecase returns ErrUserAlreadyExists", func() {
		// given
		requestBody := `{
			"email": "jhon@doe.com",
			"senha": "Abc999",
			"cpf": "72430024985",
			"nome": "Jhon Doe"
		}`
		expectedResponseBody := `{
			"message": "error creating user",
			"errors":["user already exists"]
		}`
		req, err := http.NewRequest(http.MethodPost, "/auth/users", strings.NewReader(requestBody))
		t.Require().NoError(err)
		// when
		t.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(&entity.User{}, nil)
		handlerFunc := http.HandlerFunc(t.createUserHandler.Handler)
		handlerFunc.ServeHTTP(t.rr, req)
		// then
		t.Equal(http.StatusBadRequest, t.rr.Code)
		t.JSONEq(expectedResponseBody, t.rr.Body.String())
		t.repoMock.AssertExpectations(t.T())
	})

	t.Run("must return only http status InternalServerError when usecase any other error", func() {
		// given
		requestBody := `{
			"email": "jhon@doe.com",
			"senha": "Abc999",
			"cpf": "72430024985",
			"nome": "Jhon Doe"
		}`
		req, err := http.NewRequest(http.MethodPost, "/auth/users", strings.NewReader(requestBody))
		t.Require().NoError(err)
		// when
		t.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		t.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		t.repoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.New("anything"))
		handlerFunc := http.HandlerFunc(t.createUserHandler.Handler)
		handlerFunc.ServeHTTP(t.rr, req)
		// then
		t.Equal(http.StatusInternalServerError, t.rr.Code)
		t.Empty(t.rr.Body.String())
		t.repoMock.AssertExpectations(t.T())
	})

	t.Run("must return the inserted user", func() {
		// given
		nascimento := "1996-03-06"
		cpf, name, email, senha := "72430024985", "Jhon Doe", "jhon@doe.com", "Abc999"
		requestBody := fmt.Sprintf(`{
			"cpf": "%s",
			"nome": "%s",
			"nascimento": "%s",
			"email": "%s",
			"senha": "%s"
		}`, cpf, name, nascimento, email, senha)
		req, err := http.NewRequest(http.MethodPost, "/auth/users", strings.NewReader(requestBody))
		t.Require().NoError(err)
		// when
		t.repoMock.EXPECT().GetUserByDocument(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		t.repoMock.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, database.ErrUserNotFound)
		t.repoMock.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
		handlerFunc := http.HandlerFunc(t.createUserHandler.Handler)
		handlerFunc.ServeHTTP(t.rr, req)
		// then
		t.Equal(http.StatusCreated, t.rr.Code)
		response := new(usecase.CreatedUserOutputDto)
		t.Require().NoError(json.NewDecoder(t.rr.Body).Decode(response))
		t.Equal(t.fakeId.String(), response.ID.String())
		t.Equal(cpf, response.Document)
		t.Equal(name, response.Name)
		t.Equal(nascimento, response.BirthDate.String())
		t.Equal(email, response.Email)
		t.NoError(bcrypt.CompareHashAndPassword([]byte(response.Password), []byte(senha)))
		t.repoMock.AssertExpectations(t.T())
	})
}

func TestCreateUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserHandlerTestSuite))
}

package database

import (
	"context"
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/internal/domain/entity/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

var (
	staticUUID    = uuid.New()
	validDocument = "72430024985"
)

type UserRepositoryTestSuite struct {
	repository *UserRepository
	pgxMock    pgxmock.PgxPoolIface
	suite.Suite
}

func (t *UserRepositoryTestSuite) SetupSubTest() {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	t.Require().NoError(err)
	t.pgxMock = mock
	t.repository = &UserRepository{q: New(mock)}
}

func (t *UserRepositoryTestSuite) TearDownSubTest() {
	t.pgxMock.Close()
}

func (t *UserRepositoryTestSuite) TestGetUserByDocument() {
	t.Run("must return an ErrUserNotFound when the Querier returns pgx.ErrTxClosed", func() {
		// given
		document := "12345"
		// when
		t.pgxMock.ExpectQuery(getUserByDocument).WithArgs(document).WillReturnError(pgx.ErrNoRows)
		result, err := t.repository.GetUserByDocument(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.Nil(result)
		t.Error(err)
		t.ErrorIs(err, ErrUserNotFound)
	})

	t.Run("must return an error when the Querier returns any error", func() {
		// given
		fakeErr := pgx.ErrTxClosed
		document := "12345"
		// when
		t.pgxMock.ExpectQuery(getUserByDocument).WithArgs(document).WillReturnError(fakeErr)
		result, err := t.repository.GetUserByDocument(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.Nil(result)
		t.Error(err)
		t.ErrorIs(err, fakeErr)
	})

	t.Run("must return the user from db", func() {
		// given
		document, name, email, password, birthDate := validDocument, "Jhon", "jhon@dow.com", "12345", time.Now()
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		rows := t.pgxMock.NewRows([]string{"id", "document", "name", "email", "password", "birthdate", "type"}).
			AddRow(staticUUID, document, name, email, string(hashedPass), birthDate, string(user.DefaultType))
		// when
		t.pgxMock.ExpectQuery(getUserByDocument).WithArgs(document).WillReturnRows(rows)
		result, err := t.repository.GetUserByDocument(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.NoError(err)
		t.NotNil(result)
		t.Equal(staticUUID.String(), result.ID().String())
		t.Equal(name, result.Name())
		t.Equal(email, result.Email())
		t.NoError(bcrypt.CompareHashAndPassword([]byte(result.Password()), []byte(password)))
		t.Equal(birthDate, result.BirthDate())
		t.Equal(user.DefaultType, result.Type())
	})
}

func (t *UserRepositoryTestSuite) TestGetUserByEmail() {
	t.Run("must return an ErrUserNotFound when the Querier returns pgx.ErrTxClosed", func() {
		// given
		document := "12345"
		// when
		t.pgxMock.ExpectQuery(getUserByEmail).WithArgs(document).WillReturnError(pgx.ErrNoRows)
		result, err := t.repository.GetUserByEmail(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.Nil(result)
		t.Error(err)
		t.ErrorIs(err, ErrUserNotFound)
	})

	t.Run("must return an error when the Querier returns any error", func() {
		// given
		fakeErr := pgx.ErrTxClosed
		document := "12345"
		// when
		t.pgxMock.ExpectQuery(getUserByEmail).WithArgs(document).WillReturnError(fakeErr)
		result, err := t.repository.GetUserByEmail(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.Nil(result)
		t.Error(err)
		t.ErrorIs(err, fakeErr)
	})

	t.Run("must return the user from db", func() {
		// given
		document, name, email, password, birthDate := validDocument, "Jhon", "jhon@dow.com", "12345", time.Now()
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		rows := t.pgxMock.NewRows([]string{"id", "document", "name", "email", "password", "birthdate", "type"}).
			AddRow(staticUUID, document, name, email, string(hashedPass), birthDate, string(user.DefaultType))
		// when
		t.pgxMock.ExpectQuery(getUserByEmail).WithArgs(document).WillReturnRows(rows)
		result, err := t.repository.GetUserByEmail(context.TODO(), document)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.NoError(err)
		t.NotNil(result)
		t.Equal(staticUUID.String(), result.ID().String())
		t.Equal(name, result.Name())
		t.Equal(email, result.Email())
		t.NoError(bcrypt.CompareHashAndPassword([]byte(result.Password()), []byte(password)))
		t.Equal(birthDate, result.BirthDate())
		t.Equal(user.DefaultType, result.Type())
	})
}

func (t *UserRepositoryTestSuite) TestCreate() {
	t.Run("must insert a new user", func() {
		// given
		u := new(user.User)
		// when
		t.pgxMock.ExpectExec(insertUser).
			WithArgs(u.ID(), u.Document(), u.Name(), u.Email(), u.Password(), u.BirthDate(), AuthUserType(u.Type())).
			WillReturnResult(pgxmock.NewResult("INSERT 1", 1))
		err := t.repository.Create(context.TODO(), u)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.NoError(err)
	})

	t.Run("must return an error", func() {
		// given
		u := new(user.User)
		someErr := errors.New("some error")
		// when
		t.pgxMock.ExpectExec(insertUser).
			WithArgs(u.ID(), u.Document(), u.Name(), u.Email(), u.Password(), u.BirthDate(), AuthUserType(u.Type())).
			WillReturnError(someErr)
		err := t.repository.Create(context.TODO(), u)
		// then
		t.Require().NoErrorf(t.pgxMock.ExpectationsWereMet(), "there were unfulfilled expectations: %s", err)
		t.ErrorIs(err, someErr)
	})
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

package entity

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

const validDoc = "37958959958"

func TestNewUserValidations(t *testing.T) {

	fakeErr := errors.New("fake err")
	FakeOption := func(u *User) error {
		return fakeErr
	}

	type args struct {
		document  string
		name      string
		email     string
		birthDate time.Time
		userOpts  []UserOption
	}
	tests := []struct {
		name           string
		args           args
		expectedErrStr string
	}{
		{"should not possible to create a user with an invalid document", args{document: "12345"}, "invalid cpf"},
		{"should not possible to create a user with name length less than 3", args{document: validDoc, name: "aa"}, "invalid name"},
		{"should not possible to create a user with an empty email", args{document: validDoc, name: "Jhon", email: ""}, "invalid email"},
		{"should not possible to create a user with an empty password", args{document: validDoc, name: "Jhon", email: "jhon@user.com", userOpts: []UserOption{WithPassword("")}}, "invalid password"},
		{"should not possible to create a user with a password length greater than 72", args{document: validDoc, name: "Jhon", email: "jhon@user.com", userOpts: []UserOption{WithPassword(string(make([]byte, 73)))}}, bcrypt.ErrPasswordTooLong.Error()},
		{"should not possible to create a user with when an option returns error", args{document: validDoc, name: "Jhon", email: "jhon@user.com", userOpts: []UserOption{FakeOption, WithPassword("12345")}}, fakeErr.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(uuid.New(), tt.args.document, tt.args.name, tt.args.email, tt.args.birthDate, tt.args.userOpts...)
			assert.Error(t, err)
			assert.EqualError(t, err, tt.expectedErrStr)
			assert.Nil(t, got)
		})
	}
}

func TestNewUser(t *testing.T) {
	t.Run("should create a DefaultType user by default", func(t *testing.T) {
		// given
		document, name, email, birthDate := validDoc, "Jhon", "jhon@user.com", time.Now()
		// when
		got, err := NewUser(uuid.New(), document, name, email, birthDate)
		// then
		assert.NoError(t, err)
		assert.Equal(t, DefaultType, got.Type())
		assert.Equal(t, document, got.Document())
		assert.Equal(t, name, got.Name())
		assert.Equal(t, email, got.Email())
		assert.Equal(t, birthDate, got.BirthDate())
	})

	t.Run("should create an AdminType user when the domain email is admsDomain", func(t *testing.T) {
		// given
		email := "jhon" + admsDomain
		id, document, name, birthDate := uuid.New(), validDoc, "Jhon", time.Now()
		// when
		got, err := NewUser(id, document, name, email, birthDate)
		// then
		assert.NoError(t, err)
		assert.Equal(t, AdminType, got.Type())
		assert.Equal(t, document, got.Document())
		assert.Equal(t, name, got.Name())
		assert.Equal(t, email, got.Email())
		assert.Equal(t, birthDate, got.BirthDate())
	})

	t.Run("must create a user with the provided id", func(t *testing.T) {
		// given
		id, document, name, email, birthDate := uuid.New(), validDoc, "Jhon", "jhon@doe.com", time.Now()
		// when
		got, err := NewUser(id, document, name, email, birthDate)
		// then
		assert.NoError(t, err)
		assert.Equal(t, id, got.ID())
	})

	t.Run("must create a user with the type entered in the WithType option", func(t *testing.T) {
		// given
		document, name, email, birthDate := validDoc, "Jhon", "jhon@doe.com", time.Now()
		// when
		got, err := NewUser(uuid.UUID{}, document, name, email, birthDate, WithType(AdminType))
		// then
		assert.NoError(t, err)
		assert.Equal(t, AdminType, got.Type())
	})
}

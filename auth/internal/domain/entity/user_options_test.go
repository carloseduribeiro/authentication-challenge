package entity

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestWithType(t *testing.T) {
	// given
	user := new(User)
	// when
	f := WithType(AdminType)
	err := f(user)
	// then
	assert.NoError(t, err)
	assert.Equal(t, AdminType, user.userType)
}

func TestWithPassword(t *testing.T) {
	t.Run("should not possible to create a user with an empty password", func(t *testing.T) {
		// given
		user := new(User)
		// when
		f := WithPassword("")
		err := f(user)
		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid password")
	})

	t.Run("should not possible to create a user with a password length greater than 72", func(t *testing.T) {
		// given
		user := new(User)
		// when
		f := WithPassword(string(make([]byte, 73)))
		err := f(user)
		// then
		assert.Error(t, err)
		assert.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())
	})

	t.Run("must generate the password with bcrypt DefaultCost and set it to user", func(t *testing.T) {
		// given
		user := new(User)
		password := "password"
		// when
		f := WithPassword(password)
		err := f(user)
		// then
		assert.NoError(t, err)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password)))
	})
}

func TestWithPasswordHashed(t *testing.T) {
	t.Run("should not possible to create a user with an empty hashed password", func(t *testing.T) {
		// given
		user := new(User)
		// when
		f := WithPasswordHashed("")
		err := f(user)
		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid hashed password")
	})

	t.Run("must return an error when the given password is not hashed", func(t *testing.T) {
		// given
		user := new(User)
		password := "password"
		// when
		f := WithPasswordHashed(password)
		err := f(user)
		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, bcrypt.ErrHashTooShort)
	})

	t.Run("must set user password with the given passwordHashed", func(t *testing.T) {
		// given
		user := new(User)
		password := "password"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// when
		f := WithPasswordHashed(string(hashed))
		err := f(user)
		// then
		assert.NoError(t, err)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password)))

	})
}

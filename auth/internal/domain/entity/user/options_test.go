package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithUUIDGeneratorFunc(t *testing.T) {
	t.Run("should generate and set the user id", func(t *testing.T) {
		// given
		user := new(User)
		uuidMock := uuid.New()
		uuidGeneratorFuncMock := func() (uuid.UUID, error) {
			return uuidMock, nil
		}
		// when
		f := WithUUIDGeneratorFunc(uuidGeneratorFuncMock)
		err := f(user)
		// then
		assert.NoError(t, err)
		assert.Equal(t, uuidMock, user.id)
	})

	t.Run("should return an error when the generator returns error", func(t *testing.T) {
		// given
		user := new(User)
		fakeErr := errors.New("fake err")
		uuidGeneratorFuncMock := func() (uuid.UUID, error) {
			return uuid.UUID{}, fakeErr
		}
		// when
		f := WithUUIDGeneratorFunc(uuidGeneratorFuncMock)
		err := f(user)
		// then
		assert.Error(t, err)
		assert.ErrorIs(t, err, fakeErr)
	})
}

func TestWithID(t *testing.T) {
	// given
	user := new(User)
	uuidMock := uuid.New()
	// when
	f := WithID(uuidMock)
	err := f(user)
	// then
	assert.NoError(t, err)
	assert.Equal(t, uuidMock, user.id)
}

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

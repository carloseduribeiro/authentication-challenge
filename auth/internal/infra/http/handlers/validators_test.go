package handlers

import (
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/application/usecase"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePayload(t *testing.T) {
	t.Run("must return all validation errors", func(t *testing.T) {
		// given
		args := &usecase.CreateUserInputDto{Name: "1"}
		expectedErrors := []string{
			"invalid cpf",
			"name must be at least 3 characters",
			"name must contains only letters",
			"invalid email address",
			"password must contains at least one lower case character",
			"password must contains at least one upper case character",
			"password must contains at least one numeric character",
			"password must contains at least 6 or up to 12 characters",
		}
		// then
		assert.Equal(t, expectedErrors, validateCreateUserPayload(args))
	})
}

package handlers

import (
	"github.com/carloseduribeiro/authentication-challenge/auth/internal/application/usecase"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/cpf"
	"regexp"
	"strings"
)

// validateCreateUserPayload validate the input fields and returns a slice with validation error messages
func validateCreateUserPayload(input *usecase.CreateUserInputDto) []string {
	errS := make([]string, 0, 5)
	if !cpf.Validate(input.Document) {
		errS = append(errS, "invalid cpf")
	}
	if len(strings.Trim(input.Name, " ")) < 3 {
		errS = append(errS, "name must be at least 3 characters")
	}
	if !isAlphabetical(input.Name) {
		errS = append(errS, "name must contains only letters")
	}
	if !isValidEmail(input.Email) {
		errS = append(errS, "invalid email address")
	}
	if passwdErrS := validatePasswd(input.Password); len(passwdErrS) > 0 {
		errS = append(errS, passwdErrS...)
	}
	return errS
}

func isAlphabetical(input string) bool {
	return regexp.MustCompile("^[a-zA-Z ]*$").MatchString(input)
}

func isValidEmail(input string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$").MatchString(input)
}

func validatePasswd(input string) []string {
	errS := make([]string, 0)
	if !regexp.MustCompile("[a-z]").MatchString(input) {
		errS = append(errS, "password must contains at least one lower case character")
	}
	if !regexp.MustCompile("[A-Z]").MatchString(input) {
		errS = append(errS, "password must contains at least one upper case character")
	}
	if !regexp.MustCompile("[0-9]").MatchString(input) {
		errS = append(errS, "password must contains at least one numeric character")
	}
	if !regexp.MustCompile("[a-zA-Z0-9]{6,12}").MatchString(input) {
		errS = append(errS, "password must contains at least 6 or up to 12 characters")
	}
	return errS
}

package entity

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserOption func(*User) error

func WithType(userType Type) UserOption {
	return func(u *User) error {
		u.userType = userType
		return nil
	}
}

func WithPassword(password string) UserOption {
	return func(u *User) error {
		if len(password) == 0 {
			return errors.New("invalid password")
		}
		pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.password = string(pass)
		return nil
	}
}

func WithPasswordHashed(passwordHashed string) UserOption {
	return func(u *User) error {
		if len(passwordHashed) == 0 {
			return errors.New("invalid hashed password")
		}
		if _, err := bcrypt.Cost([]byte(passwordHashed)); err != nil {
			return err
		}
		u.password = passwordHashed
		return nil
	}
}

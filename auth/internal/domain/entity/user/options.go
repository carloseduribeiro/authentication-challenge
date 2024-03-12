package user

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Option func(*User) error

func WithUUIDGeneratorFunc(uuidGeneratorFunc func() (uuid.UUID, error)) Option {
	return func(u *User) error {
		id, err := uuidGeneratorFunc()
		if err != nil {
			return err
		}
		u.id = id
		return nil
	}
}

func WithID(id uuid.UUID) Option {
	return func(u *User) error {
		u.id = id
		return nil
	}
}

func WithType(userType Type) Option {
	return func(u *User) error {
		u.userType = userType
		return nil
	}
}

func WithPassword(password string) Option {
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

func WithPasswordHashed(passwordHashed string) Option {
	return func(u *User) error {
		if len(passwordHashed) == 0 {
			return errors.New("invalid password")
		}
		if _, err := bcrypt.Cost([]byte(passwordHashed)); err != nil {
			return err
		}
		u.password = passwordHashed
		return nil
	}
}

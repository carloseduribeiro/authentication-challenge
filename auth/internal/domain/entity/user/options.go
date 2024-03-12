package user

import (
	"github.com/google/uuid"
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

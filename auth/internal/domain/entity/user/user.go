package user

import (
	"errors"
	"github.com/carloseduribeiro/auth-challenge/auth/pkg/domain/entity/cpf"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const admsDomain = "@br.furabolso.com"

type User struct {
	id        uuid.UUID
	userType  Type
	document  string
	name      string
	email     string
	password  string
	birthDate time.Time
}

// NewUser creates a new user.
func NewUser(document, name, email, password string, birthDate time.Time, userOpts ...Option) (*User, error) {
	if !cpf.Validate(document) {
		return nil, errors.New("invalid cpf")
	}
	if len(name) < 3 {
		return nil, errors.New("invalid name")
	}
	if len(email) < 3 {
		return nil, errors.New("invalid email")
	}
	if len(password) == 0 {
		return nil, errors.New("invalid password")
	}
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	user := &User{
		id:        id,
		userType:  DefaultType,
		document:  document,
		name:      name,
		email:     email,
		password:  string(pass),
		birthDate: birthDate,
	}
	if strings.Contains(email, admsDomain) {
		user.userType = AdminType
	}
	for _, opt := range userOpts {
		if err = opt(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Type() Type {
	return u.userType
}

func (u *User) Document() string {
	return u.document
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) BirthDate() time.Time {
	return u.birthDate
}

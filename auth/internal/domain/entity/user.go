package entity

import (
	"errors"
	"fmt"
	"github.com/carloseduribeiro/authentication-challenge/lib-utils/pkg/cpf"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Type string

const (
	DefaultType Type = "default"
	AdminType   Type = "admin"
)

const admsDomain = "@br.furabolso.com"

type User struct {
	id            uuid.UUID
	userType      Type
	document      string
	name          string
	email         string
	password      string
	birthDate     time.Time
	activeSession *Session
}

// NewUser creates a new user
func NewUser(id uuid.UUID, document, name, email string, birthDate time.Time, userOpts ...UserOption) (*User, error) {
	if !cpf.Validate(document) {
		return nil, errors.New("invalid cpf")
	}
	if len(name) < 3 {
		return nil, errors.New("invalid name")
	}
	if len(email) < 3 {
		return nil, errors.New("invalid email")
	}
	user := &User{
		id:        id,
		userType:  DefaultType,
		document:  document,
		name:      name,
		email:     email,
		birthDate: birthDate,
	}
	if strings.Contains(email, admsDomain) {
		user.userType = AdminType
	}
	for _, opt := range userOpts {
		if err := opt(user); err != nil {
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

func (u *User) NewSession(sessionId uuid.UUID, createdAt time.Time, d time.Duration, secretKey string) (string, error) {
	u.activeSession = NewSession(sessionId, u.id, createdAt, d)
	mapClaims := jwt.MapClaims{
		"id":       sessionId,
		"userId":   u.id,
		"userType": u.userType,
		"exp":      u.activeSession.ExpiresAt(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing jwt: %v", err)
	}
	return tokenStr, nil
}

func (u *User) ActiveSession() *Session {
	return u.activeSession
}

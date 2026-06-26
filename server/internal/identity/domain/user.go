package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	id          uuid.UUID
	username    Username
	phoneNumber PhoneNumber
	credential  UserCredential
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   time.Time
}

func NewUser(username, phoneNumber string) (*User, error) {
	name, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        uuid.New(),
		username:  name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (u *User) UpdatePhoneNumber(phoneNumber string) error {
	pn, err := NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	u.phoneNumber = pn
	u.updatedAt = time.Now()
	return nil
}

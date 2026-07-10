package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Username    Username
	Email       Email
	PhoneNumber PhoneNumber
	Credential  UserCredential
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func NewUser(username Username) (*User, error) {
	return &User{
		ID:        uuid.New(),
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) UpdatePhoneNumber(phoneNumber string) error {
	pn, err := NewPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	u.PhoneNumber = pn
	u.UpdatedAt = time.Now()
	return nil
}

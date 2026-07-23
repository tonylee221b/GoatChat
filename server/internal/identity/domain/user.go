package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username Username
	Contact  Contact
	Audit    Audit
}

func NewUser(username Username) *User {
	return &User{
		ID:       uuid.New(),
		Username: username,
	}
}

func (u *User) UpdateContact(contact Contact) error {
	c, err := NewContact(contact.PhoneNumber.Value, contact.Email.Value)
	if err != nil {
		return err
	}

	u.Contact = c
	return nil
}

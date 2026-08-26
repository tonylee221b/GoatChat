package domain

import (
	"github.com/google/uuid"
)

// User struct is an Aggregate Root
//
// it manages user data
type User struct {
	ID           uuid.UUID
	Username     Username
	PasswordHash PasswordHash
	Contact      Contact
	Audit        Audit
}

func NewUser(username Username, pwHash PasswordHash) *User {
	return &User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: pwHash,
	}
}

func RestoreUser(
	id uuid.UUID,
	username Username,
	pwHash PasswordHash,
	contact Contact,
	audit Audit,
) *User {
	return &User{
		ID:           id,
		Username:     username,
		PasswordHash: pwHash,
		Contact:      contact,
		Audit:        audit,
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

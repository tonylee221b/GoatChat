package domain

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

type UserStatus string

const (
	ACTIVE   UserStatus = "active"
	INACTIVE UserStatus = "inactive"
)

type Username struct {
	Value string
}

func NewUsername(name string) (Username, error) {
	if name == "" {
		slog.Info("username is empty", "username", name)
		return Username{}, errors.New(ErrInvalidUsername)
	}

	return Username{Value: name}, nil
}

type Contact struct {
	PhoneNumber
	Email
}

func NewContact(phoneNumber, email string) (Contact, error) {
	pn, err := newPhoneNumber(phoneNumber)
	if err != nil {
		return Contact{}, err
	}

	e, err := newEmail(email)
	if err != nil {
		return Contact{}, err
	}

	return Contact{
		PhoneNumber: pn,
		Email:       e,
	}, nil
}

type PhoneNumber struct {
	Value string
}

func newPhoneNumber(number string) (PhoneNumber, error) {
	if number == "" {
		slog.Info("phone number is empty", "phone number", number)
		return PhoneNumber{}, errors.New(ErrInvalidPhoneNumberFormat)
	}

	normalized := strings.NewReplacer(
		"-", "",
		" ", "",
	).Replace(number)

	koreanPhonePattern := regexp.MustCompile(`^(?:010[0-9]{4}[0-9]{4}|01[16789][0-9]{3,4}[0-9]{4})$`)
	if ok := koreanPhonePattern.MatchString(normalized); !ok {
		slog.Info("invalid korean phone number", "phone number", koreanPhonePattern)
		return PhoneNumber{}, errors.New(ErrInvalidPhoneNumberFormat)
	}

	return PhoneNumber{Value: number}, nil
}

type Email struct {
	Value string
}

func newEmail(email string) (Email, error) {
	if email == "" {
		slog.Info("email is empty", "email", email)
		return Email{}, errors.New(ErrInvalidEmailFormat)
	}

	return Email{Value: email}, nil
}

type Audit struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewAudit(createdAt, updatedAt time.Time, deletedAt *time.Time) Audit {
	return Audit{createdAt, updatedAt, deletedAt}
}

type PasswordHash struct {
	Value string
}

func NewPasswordHash(hashedPW string) (PasswordHash, error) {
	if strings.TrimSpace(hashedPW) == "" {
		slog.Info("password is blank")
		return PasswordHash{}, errors.New(ErrInvalidPasswordFormat)
	}

	return PasswordHash{
		Value: hashedPW,
	}, nil
}

package domain

import (
	"errors"
	"regexp"
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
		return Username{}, errors.New(ErrInvalidUsername)
	}

	return Username{Value: name}, nil
}

type PhoneNumber struct {
	Value string
}

func NewPhoneNumber(number string) (PhoneNumber, error) {
	if number == "" {
		return PhoneNumber{}, errors.New(ErrInvalidPhoneNumberFormat)
	}

	if ok, _ := regexp.MatchString("", number); !ok {
		return PhoneNumber{}, errors.New(ErrInvalidPhoneNumberFormat)
	}

	return PhoneNumber{Value: number}, nil
}

type Email struct {
	Value string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New(ErrInvalidEmailFormat)
	}

	return Email{Value: email}, nil
}

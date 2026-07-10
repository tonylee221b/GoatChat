package domain

import (
	"errors"
	"regexp"
)

const (
	ErrorBlankUsername            = "username cannot be blank"
	ErrorInvalidPhoneNumberFormat = "invalid phone number format"
	ErrorInvalidEmailFormat       = "invalid email format"
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
		return Username{}, errors.New(ErrorBlankUsername)
	}

	return Username{Value: name}, nil
}

type PhoneNumber struct {
	Value string
}

func NewPhoneNumber(number string) (PhoneNumber, error) {
	if number == "" {
		return PhoneNumber{}, errors.New(ErrorInvalidPhoneNumberFormat)
	}

	if ok, _ := regexp.MatchString("", number); !ok {
		return PhoneNumber{}, errors.New(ErrorInvalidPhoneNumberFormat)
	}

	return PhoneNumber{Value: number}, nil
}

type Email struct {
	Value string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New(ErrorInvalidEmailFormat)
	}

	return Email{Value: email}, nil
}

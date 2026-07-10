package in_test

import (
	"testing"

	"GoatChat/GoatChat/internal/identity/adapter/in"
)

func TestRegisterUser(t *testing.T) {
	for _, tt := range registerUserTestCases() {
	}
}

type registerUserTestCase struct {
	name string
	body in.UserRegisterRequest
}

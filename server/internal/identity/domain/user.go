package domain

import "github.com/google/uuid"

type User struct {
	id         uuid.UUID
	credential UserCredential
}

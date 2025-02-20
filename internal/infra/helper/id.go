package helper

import "github.com/google/uuid"

type ID = uuid.UUID

// VO - Value Object
func NewID() ID {
	return uuid.New()
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

package utils

import "github.com/google/uuid"

func GenerateUuid() string {
	newUuid := uuid.New()
	uuid := newUuid.String()
	return uuid
}

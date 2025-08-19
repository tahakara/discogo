package utils

import (
	"time"

	"github.com/google/uuid"
)

func ValidateUUID(id string) (bool, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetFormatedCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

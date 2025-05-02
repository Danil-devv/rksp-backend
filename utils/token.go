package utils

import (
	"errors"
	"os"
)

func GetSecretPair() (string, string, error) {
	key := os.Getenv("TURN_KEY")
	value := os.Getenv("TURN_VALUE")

	if key == "" || value == "" {
		return "", "", errors.New("TURN_KEY or TURN_VALUE is empty")
	}

	return key, value, nil
}

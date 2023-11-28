package utils

import (
	"errors"
	"math/rand"
	"time"
	"unicode"
)

func CheckGateway(name string) error {
	err := hasDigits(name)
	if err != nil {
		return err
	}
	return nil
}
func hasDigits(name string) error {
	hasDigits := false

	for _, char := range name {
		if unicode.IsDigit(char) {
			hasDigits = true
			break
		}
	}

	if hasDigits {
		return errors.New("the gateway name must be string")
	}
	return nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

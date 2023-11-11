package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Generate a salt and hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	// Compare the given password with the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

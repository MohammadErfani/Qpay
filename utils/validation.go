package utils

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func IsValidPhoneNumber(phoneNumber string) bool {
	regexPattern := `^\d{9}$`
	phoneNumberRegex := regexp.MustCompile(regexPattern)
	return phoneNumberRegex.MatchString(phoneNumber)
}

func IsValidNationalCode(nationalCode string) bool {
	regexPattern := `^\d{10}$`
	nationalCodeRegex := regexp.MustCompile(regexPattern)
	return nationalCodeRegex.MatchString(nationalCode)
}

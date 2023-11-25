package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func IsValidEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email is not in correct format")
	}
	return nil
}
func IsValidPhoneNumber(phoneNumber string) error {
	// regexPattern := `^\d{9}$`
	regexPattern := `^(\+98|0)?9\d{9}$`
	phoneNumberRegex := regexp.MustCompile(regexPattern)
	if !phoneNumberRegex.MatchString(phoneNumber) {
		return errors.New("phone number is not in correct format")
	}
	return nil
}

func IsValidNationalCode(nationalCode string) error {
	regexPattern := `^\d{10}$`
	nationalCodeRegex := regexp.MustCompile(regexPattern)
	if !nationalCodeRegex.MatchString(nationalCode) {
		return errors.New("identity is not in correct format")
	}
	return nil
}

func IsRequired(requiredFields map[string]string) error {
	for fieldName, value := range requiredFields {
		if len(strings.TrimSpace(value)) == 0 {
			return errors.New(fmt.Sprintf("%s is required", fieldName))
		}
	}
	return nil
}

func IsRequiredID(requiredFields map[string]uint) error {
	for fieldName, value := range requiredFields {
		if value == 0 {
			return errors.New(fmt.Sprintf("%v is required", fieldName))
		}
	}
	return nil
}

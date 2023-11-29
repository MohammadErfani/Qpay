package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidPhoneNumber(t *testing.T) {
	// correct phone number
	assert.Nil(t, IsValidPhoneNumber("09122161259"))
	assert.Nil(t, IsValidPhoneNumber("+989122161259"))
	// test length
	assert.Error(t, IsValidPhoneNumber("091221659"))
	// test for invalid phone number entry
	assert.Error(t, IsValidPhoneNumber("0912216m259"))
	assert.Error(t, IsValidPhoneNumber("0912216250_"))
	assert.Error(t, IsValidPhoneNumber("49122162500"))
	assert.Error(t, IsValidPhoneNumber("01122162500"))
}

func TestIsValidEmail(t *testing.T) {
	assert.Nil(t, IsValidEmail("mohammad@gmail.com"))
	assert.Nil(t, IsValidEmail("admin@quera.com"))
	assert.Nil(t, IsValidEmail("123_mohammad.erfani@quera.com"))
	assert.Error(t, IsValidEmail("mohammad_gmail.com"))
	assert.Error(t, IsValidEmail("mohammad@.com"))
	assert.Error(t, IsValidEmail("mohammad@gmail."))
	assert.Error(t, IsValidEmail("mohammad"))
}

func TestIsValidNationalCode(t *testing.T) {
	assert.Nil(t, IsValidNationalCode("0441111111"))
	assert.Nil(t, IsValidNationalCode("1234567890"))
	assert.Error(t, IsValidNationalCode("044111111"))
	assert.Error(t, IsValidNationalCode("044111mm11"))
	assert.Error(t, IsValidNationalCode("044111@111"))

}

func TestIsRequired(t *testing.T) {
	required := map[string]string{
		"email":    "mohammad@gmail.com",
		"password": "12345",
	}
	assert.Nil(t, IsRequired(required))
	required["email"] = ""
	assert.Error(t, IsRequired(required))
	required["email"] = " "
	assert.Error(t, IsRequired(required))

}

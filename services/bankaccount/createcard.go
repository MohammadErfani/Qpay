package bankaccount

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

func CheckSheba(sheba string) error {
	err := checkShebaCorrect(sheba)
	if err != nil {
		return err
	}
	return nil
}

func checkShebaCorrect(sheba string) error {
	LenSheba := len(sheba)
	if LenSheba != 26 {
		return errors.New("the length of sheba code is not correct")
	}
	if !strings.HasPrefix(sheba, "IR") {
		return errors.New("the prefix of sheba code is not correct")
	}
	switch sheba[:5] {
	case "IR740":
		return nil
	case "IR020":
		return nil
	case "IR510":
		return nil
	case "IR410":
		return nil
	default:
		return errors.New("this card Is not supported")
	}
}

func GetOwnerSheba(sheba string) (uint, string) {
	var bankID uint
	switch sheba[:5] {
	case "IR740":
		bankID = 1
	case "IR020":
		bankID = 2
	case "IR510":
		bankID = 3
	case "IR410":
		bankID = 4
	}

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(7)
	var owner string
	switch randomNumber {
	case 0:
		owner = "علی عباسی"
	case 1:
		owner = "محمد عرفانی"
	case 2:
		owner = "امین ظهرابی"
	case 4:
		owner = "مسعود اقدسی فام"
	case 5:
		owner = "امید مقدس"
	case 6:
		owner = "ایلیا میرزائی"
	}
	return bankID, owner

	//IR740170000000106748249001
}

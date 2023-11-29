package utils

import (
	"errors"
	"math/rand"
	"time"
)

var bankCode = map[string]string{
	"1000": "ملی ایران",
	"1011": "آینده",
	"1022": "اقتصاد نوین",
	"1033": "ایران زمین",
	"1044": "پارسیان",
	"1055": "پاسارگاد",
	"1066": "تجارت",
	"1077": "سپه",
	"1088": "توسعه تعاون",
	"1099": "کشاورزی",
}

func GetIdentityAndBank(sheba string) (identity, bankName string) {
	if len(sheba) != 24 {
		return "", ""
	}
	bankName = bankCode[sheba[:4]]
	identity = sheba[4:14]
	return identity, bankName
}

func CheckSheba(sheba string) error {
	return checkShebaCorrect(sheba)
}

func checkShebaCorrect(sheba string) error {
	LenSheba := len(sheba)
	if LenSheba != 24 {
		return errors.New("the length of sheba code is not correct")
	}
	return nil
	//if !strings.HasPrefix(sheba, "IR") {
	//	return errors.New("the prefix of sheba code is not correct")
	//}
	//switch sheba[:5] {
	//case "IR740":
	//	return nil
	//case "IR020":
	//	return nil
	//case "IR510":
	//	return nil
	//case "IR410":
	//	return nil
	//default:
	//	return errors.New("this card Is not supported")
	//}
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

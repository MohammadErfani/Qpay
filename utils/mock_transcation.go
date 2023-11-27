package utils

import (
	"Qpay/models"
	"errors"
	"fmt"
)

type PurchaserCard struct {
	CardNumber string
	CardYear   int
	CardMonth  int
	CVV2       int
	Password   int
	Owner      string
}

func Transaction(paymentAmount float64, cardYear, cardMonth int, cardNumber string, cvv2, password int) (string, error) {
	//if err := cardExpireCheck(CardYear, CardMonth); err != nil {
	//	return "", errors.New("card is expire")
	//}
	//if err := IsValidPhoneNumber(PhoneNumber); err != nil {
	//	return "", errors.New("phone number is not correct")
	//}
	card := PurchaserCard{
		CardNumber: cardNumber,
		CardMonth:  cardMonth,
		CardYear:   cardYear,
		CVV2:       cvv2,
		Password:   password,
	}
	if _, err := CheckCardCredential(card); err != nil {
		return "", err
	}
	if paymentAmount < 10000 || paymentAmount > 1200000000 {
		return "", errors.New("this payment is not correct")
	}
	return GenerateTrackingCode(models.Paid), nil
}
func PurchaserBankAccount(PurchaserCard string) string {
	switch PurchaserCard {
	case "6037991123787766":
		return "mohammad erfani"
	case "5042061042417820":
		return "amin zohrabi"
	case "5042061042417830":
		return "kiana hoseini"
	case "1234567812345678":
		return "omid moghadas"
	default:
		return ""
	}
}

//func cardExpireCheck(CardYear int, CardMonth int) error {
//	if CardYear < 1402 {
//		return errors.New("card is expire")
//	}
//	if CardYear == 1402 {
//		if CardMonth >= 9 {
//			return nil
//		}
//		if CardMonth < 9 {
//			return errors.New("card is expire")
//		}
//		return errors.New("card is expire")
//	}
//	if CardYear > 1402 {
//		return nil
//	}
//	return nil
//}

func GenerateTrackingCode(Status uint8) string {
	return fmt.Sprintf("%v%v%v",
		GenerateRandomString(4),
		Status,
		GenerateRandomString(5))
}

func CheckCardCredential(card PurchaserCard) (string, error) {
	cards := []PurchaserCard{
		{
			CardNumber: "6037991123787766",
			CardMonth:  2,
			CardYear:   1403,
			CVV2:       111,
			Password:   123456,
			Owner:      "Mohammad Erfani",
		},
		{
			CardNumber: "5042061042417820",
			CardMonth:  3,
			CardYear:   1404,
			CVV2:       222,
			Password:   123123,
			Owner:      "Amin Zohrabi",
		},
		{
			CardNumber: "5042061042417830",
			CardMonth:  4,
			CardYear:   1405,
			CVV2:       333,
			Password:   654321,
			Owner:      "Kiana Hoseini",
		},
		{
			CardNumber: "1234567812345678",
			CardMonth:  5,
			CardYear:   1406,
			CVV2:       444,
			Password:   101010,
			Owner:      "Omid Moghadas",
		},
	}
	for _, pc := range cards {
		if card.CardNumber == pc.CardNumber {
			if card.CardMonth == pc.CardMonth && card.CardYear == pc.CardYear && card.Password == pc.Password && card.CVV2 == pc.CVV2 {
				return pc.Owner, nil
			}
			return "", errors.New("credential doesn't match")
		}
	}
	return "", errors.New("card is incorrect")

}

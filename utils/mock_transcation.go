package utils

import (
	"errors"
	"slices"
)

func Transaction(PaymentAmount float64, CardYear int, CardMonth int, PhoneNumber, PurchaserCard string) error {
	if err := cardExpireCheck(CardYear, CardMonth); err != nil {
		return errors.New("card is expire")
	}
	if err := IsValidPhoneNumber(PhoneNumber); err != nil {
		return errors.New("phone number is not correct")
	}
	rightCardNumber := []string{"6037991123787766", "5042061042417820", "5042061042417830"}
	if !slices.Contains(rightCardNumber, PurchaserCard) {
		return errors.New("card number is not correct")
	}
	if PaymentAmount < 10000 || PaymentAmount > 1200000000 {
		return errors.New("this payment is not correct")
	}
	return nil
}
func PurchaserBankAccount(PurchaserCard string) string {
	switch PurchaserCard {
	case "6037991123787766":
		return "mohammad erfani"
	case "5042061042417820":
		return "amin zohrabi"
	case "5042061042417830":
		return "kiana hoseini"
	default:
		return ""
	}
}
func ComisionCalc(PaymentAmount, comision float64) float64 {
	return PaymentAmount / comision
}
func cardExpireCheck(CardYear int, CardMonth int) error {
	if CardYear < 1402 {
		return errors.New("card is expire")
	}
	if CardYear == 1402 {
		if CardMonth >= 9 {
			return nil
		}
		if CardMonth < 9 {
			return errors.New("card is expire")
		}
		return errors.New("card is expire")
	}
	if CardYear > 1402 {
		return nil
	}
	return nil
}

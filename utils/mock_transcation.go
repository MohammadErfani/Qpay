package utils

import "errors"

var cardNumber = map[string]uint{
	"6037991123787766": 1,
	"5042061042417820": 2,
	"5042061042417830": 3,
	"5042061042417840": 2,
	"5042061042417850": 2,
	"5042061042417860": 3,
	"5042061042417870": 1,
	"5042061042417880": 1,
	"5042061042417890": 1,
}

func Transaction(PaymentAmount float64, CardYear int, CardMonth int, PhoneNumber string) error {
	if err := cardExpireCheck(CardYear, CardMonth); err != nil {
		return errors.New("card is expire")
	}
	return nil

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

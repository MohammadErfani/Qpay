package bankaccount

import (
	"Qpay/models"
	"gorm.io/gorm"
)

func CreateCard(db *gorm.DB, card models.BankAccount) (*gorm.DB, error) {
	//password, _ := utils.HashPassword("1234")
	//user := models.User{
	//	Name:        "Mohammad Erfani",
	//	Email:       "mohammad@gmail.com",
	//	Username:    "mohammadErfani",
	//	Password:    password,
	//	PhoneNumber: "09121111111",
	//	Address:     "Tehran,...",
	//	Identity:    "0441111111",
	//	Role:        models.IsNaturalPerson,
	//}
	//err := db.FirstOrCreate(&user, models.User{Email: user.Email, PhoneNumber: user.PhoneNumber, Username: user.Username}).Error
	//if err != nil {
	//	log.Fatal(err)
	//}
	result := db.Limit(3).Find(&models.User{Username: "mohammadErfani"})
	return result, nil
}

package cmd

import (
	"Qpay/config"
	"Qpay/database"
	"Qpay/models"
	"Qpay/utils"
	"github.com/spf13/cobra"
	"log"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed fake data to database",
	Long: `seed fake data to database:
	fill database with fake data of users, gateways, commissions, banks,...`,
	Run: func(cmd *cobra.Command, args []string) {
		seed(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func seed(configPath string) {
	cfg := config.InitConfig(configPath)
	db := database.NewPostgres(cfg)
	password, _ := utils.HashPassword("1234")
	user := models.User{
		Name:        "Mohammad Erfani",
		Email:       "mohammad@gmail.com",
		Username:    "mohammadErfani",
		Password:    password,
		PhoneNumber: "09121111111",
		Address:     "Tehran,...",
		Identity:    "0441111111",
		Role:        models.IsNaturalPerson,
	}
	err := db.FirstOrCreate(&user, models.User{Email: user.Email, PhoneNumber: user.PhoneNumber, Username: user.Username}).Error
	if err != nil {
		log.Fatal(err)
	}

	bank := []models.Bank{
		{
			Name: "بانک ملی ایران",
			Logo: "https://bmi.ir/app_themes/faresponsive/img/bmilogo.png",
		},
		{
			Name: "بانک آینده",
			Logo: "https://static.idpay.ir/banks/ayandeh.png",
		},
		{
			Name: "بانک اقتصاد نوین",
			Logo: "https://static.idpay.ir/banks/eghtesad-novin.png",
		},
		{
			Name: "بانک ایران زمین",
			Logo: "https://static.idpay.ir/banks/iran-zamin.png",
		},
		{
			Name: "بانک پارسیان",
			Logo: "https://static.idpay.ir/banks/parsian.png",
		},
		{
			Name: "بانک پاسارگاد",
			Logo: "https://static.idpay.ir/banks/pasargad.png",
		},
		{
			Name: "بانک تجارت",
			Logo: "https://static.idpay.ir/banks/tejarat.png",
		},
		{
			Name: "بانک سپه",
			Logo: "https://static.idpay.ir/banks/sepah.png",
		},
		{
			Name: "بانک توسعه تعاون",
			Logo: "https://static.idpay.ir/banks/tosee-taavon.png",
		},
		{
			Name: "بانک کشاورزی",
			Logo: "https://static.idpay.ir/banks/keshavarzi.png",
		},
		{
			Name: "بانک مسکن",
			Logo: "https://static.idpay.ir/banks/maskan.png",
		}}
	db.Save(&bank)

}

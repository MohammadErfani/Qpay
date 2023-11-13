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

	bank := []models.Bank{{
		Name: "بانک ملی ایران",
		Logo: "https://bmi.ir/app_themes/faresponsive/img/bmilogo.png",
	},
		{
			Name: "بانک سپه",
			Logo: "https://www.banksepah.ir/UI/Styles/sepah1401/images/logo.png",
		},
		{
			Name: "بانک صنعت و معدن",
			Logo: "https://www.bim.ir/_Bimportal/logos/fa-IR/logo.png",
		},
		{
			Name: "بانک کشاورزی",
			Logo: "https://www.bki.ir/Portals/0/Logo.jpg",
		},
		{
			Name: "بانک مسکن",
			Logo: "https://www.bank-maskan.ir/image/layout_set_logo?img_id=15494931&t=1699877843772",
		}}
	db.Save(&bank)

}

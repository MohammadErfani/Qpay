package cmd

import (
	"Qpay/config"
	"Qpay/database"
	"Qpay/models"
	"Qpay/services"
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
	// seeding User
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
	err := db.Where(models.User{Email: user.Email}).Or(models.User{PhoneNumber: user.PhoneNumber}).Or(models.User{Username: user.Username}).FirstOrCreate(&user).Error
	if err != nil {
		log.Fatal("Error seeding User", err)
	}
	// Seeding Banks
	bankNames := []string{"Pasargad", "ملی", "ملت", "Eghtesad Novin", "Tejarat"}
	for _, bn := range bankNames {
		err = db.Where(models.Bank{Name: bn}).FirstOrCreate(&models.Bank{Name: bn}).Error
		if err != nil {
			log.Fatal("Error seeding Banks", err)
		}
	}

	// seeding commission
	var i float64 = 2
	for ; i <= 4; i += 2 {
		amount := 100 * i
		percent := 0.02 / i
		err = db.FirstOrCreate(&models.Commission{
			PercentagePerTrans: percent,
			AmountPerTrans:     amount,
			Status:             models.CommIsActive,
		}, models.Commission{PercentagePerTrans: percent, AmountPerTrans: amount}).Error
		if err != nil {
			log.Fatal("Error seeding commission", err)
		}
	}
	// seed bank account
	//var bank models.Bank
	//err = db.First(&bank, models.Bank{Name: "Mellat"}).Error
	//if err != nil {
	//	log.Fatal("Error getting Bank")
	//}
	bankAccount := models.BankAccount{
		Sheba: "101104411111111234123412",
	}
	err = services.SetUserAndBankForBankAccount(db, &bankAccount)
	if err != nil {
		log.Fatal("Error in Sheba: ", err)
	}
	err = db.FirstOrCreate(&bankAccount, models.BankAccount{Sheba: bankAccount.Sheba}).Error
	if err != nil {
		log.Fatal("error seeding bank account", err)
	}
	//err = db.FirstOrCreate(&bankAccount).Error
}

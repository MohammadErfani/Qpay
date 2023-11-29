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
		Name:        "Omid Moghadas",
		Email:       "electromidz@gmail.com",
		Username:    "o.moghadas",
		Password:    password,
		PhoneNumber: "09155225920",
		Address:     "Mashhad",
		Identity:    "9800",
		Role:        models.IsNaturalPerson,
	}
	err := db.Where(models.User{Email: user.Email}).Or(models.User{PhoneNumber: user.PhoneNumber}).Or(models.User{Username: user.Username}).FirstOrCreate(&user).Error
	if err != nil {
		log.Fatal("Error seeding User", err)
	}
	password, _ = utils.HashPassword("1234")
	user = models.User{
		Name:        "Mohammad Erfani",
		Email:       "mohammad@gmail.com",
		Username:    "mohammadErfani",
		Password:    password,
		PhoneNumber: "09121111111",
		Address:     "Tehran,...",
		Identity:    "0441111111",
		Role:        models.IsNaturalPerson,
	}
	err = db.Where(models.User{Email: user.Email}).Or(models.User{PhoneNumber: user.PhoneNumber}).Or(models.User{Username: user.Username}).FirstOrCreate(&user).Error
	if err != nil {
		log.Fatal("Error seeding User", err)
	}
	banks := []models.Bank{
		{
			Name: "ملی ایران",
			Logo: "https://bmi.ir/app_themes/faresponsive/img/bmilogo.png",
		},
		{
			Name: "آینده",
			Logo: "https://static.idpay.ir/banks/ayandeh.png",
		},
		{
			Name: "اقتصاد نوین",
			Logo: "https://static.idpay.ir/banks/eghtesad-novin.png",
		},
		{
			Name: "ایران زمین",
			Logo: "https://static.idpay.ir/banks/iran-zamin.png",
		},
		{
			Name: "پارسیان",
			Logo: "https://static.idpay.ir/banks/parsian.png",
		},
		{
			Name: "پاسارگاد",
			Logo: "https://static.idpay.ir/banks/pasargad.png",
		},
		{
			Name: "تجارت",
			Logo: "https://static.idpay.ir/banks/tejarat.png",
		},
		{
			Name: "سپه",
			Logo: "https://static.idpay.ir/banks/sepah.png",
		},
		{
			Name: "توسعه تعاون",
			Logo: "https://static.idpay.ir/banks/tosee-taavon.png",
		},
		{
			Name: "کشاورزی",
			Logo: "https://static.idpay.ir/banks/keshavarzi.png",
		},
		{
			Name: "مسکن",
			Logo: "https://static.idpay.ir/banks/maskan.png",
		}}

	for _, bank := range banks {
		err = db.Where(models.Bank{Name: bank.Name}).FirstOrCreate(&bank).Error
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
	err = services.SetUserAndBankForBankAccount(db, user.ID, &bankAccount)
	if err != nil {
		log.Fatal("Error in Sheba: ", err)
	}
	err = db.FirstOrCreate(&bankAccount, models.BankAccount{Sheba: bankAccount.Sheba}).Error
	if err != nil {
		log.Fatal("error seeding bank account", err)
	}

	admin := models.User{
		Name:     cfg.Admin.Name,
		Email:    cfg.Admin.Email,
		Username: cfg.Admin.Username,
		Role:     models.IsAdmin,
	}
	admin.Password, _ = utils.HashPassword(cfg.Admin.Password)
	err = db.Where("email=?", admin.Email).Or("username=?", admin.Username).FirstOrCreate(&admin).Error
	if err != nil {
		log.Fatal("error seeding admin", err)
	}
	//err = db.FirstOrCreate(&bankAccount).Error
}

package cmd

import (
	"Qpay/config"
	"Qpay/database"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/cobra"
	"log"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate tables",
	Long: `migrate tables:
	base on state and steps, migrate the database 
`,
	Run: func(cmd *cobra.Command, args []string) {
		state, _ := cmd.Flags().GetString("state")
		steps, _ := cmd.Flags().GetInt("steps")
		migratePostgres(cfgFile, state, steps)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().String("state", "up", "write the state")
	migrateCmd.Flags().Int("steps", 1, "write the steps that you need up or down")
}

func migratePostgres(configPath string, state string, steps int) {
	cfg := config.InitConfig(configPath)
	db := database.NewPostgres(cfg)
	sql, _ := db.DB()
	driver, err := postgres.WithInstance(sql, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://database/migration",
		"postgres",
		driver)
	if err != nil {
		log.Fatal(err)
	}

	switch state {
	case "up":
		err = mig.Up()
		if err != nil && err.Error() != "no change" {
			log.Fatal(err)
		}
		log.Println("migrate up has done")
	case "down":
		err = mig.Down()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migrate down has done")
	case "drop":
		err = mig.Drop()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migrate drop has done")
	case "steps":
		err = mig.Steps(steps)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("migration with steps has done")
	default:
		log.Fatal("nothing")
	}
}

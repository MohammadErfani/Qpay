package cmd

import (
	"Qpay/config"
	"Qpay/database"
	"Qpay/server"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "Qpay",
	Short: "Start program",
	Long: `Start Program:
	read config file from config.yaml
	connect to database
	start serving 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.InitConfig(cfgFile)
		db := database.NewPostgres(cfg)
		_ = db
		server.StartServer(cfg, server.Instance())
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file")
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed fake data to database",
	Long: `seed fake data to database:
	fill database with fake data of users, gateways, commissions, banks,...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seed called")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

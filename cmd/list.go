package cmd

import (
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/miflo"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List pending migrations",
	Long:    "The list command lists all migrations that have not been applied in the migrations directory.",
	Args:    cobra.NoArgs,
	Example: "miflo list",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}

		databaseConnection := os.Getenv("DATABASE_URL")
		if databaseConnection == "" {
			fmt.Println("DATABASE_URL is not set")
			return
		}

		database, err := database.NewDatabase(databaseConnection)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer database.Close()

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("error getting current working directory")
		}

		if err := miflo.ListPendingMigrations(database, cwd); err != nil {
			fmt.Println(err)
			return
		}
	},
}

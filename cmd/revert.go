package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/miflo"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(revertCmd)
}

var revertCmd = &cobra.Command{
	Use:     "revert",
	Short:   "Revert lat migration",
	Long:    "The revert command rolls back all the database migrations that were most recently applied using the up command.",
	Args:    cobra.NoArgs,
	Example: "miflo revert",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("error getting current working directory: ", err)
		}

		databaseConnection := os.Getenv("DATABASE_URL")
		if databaseConnection == "" {
			fmt.Println("DATABASE_URL is not set")
			return
		}

		database, err := database.NewDatabase(databaseConnection)
		if err != nil {
			fmt.Println("Error setting up database:", err)
			return
		}

		defer database.Close()

		ctx := context.Background()

		if err := miflo.RevertMigrations(database, ctx, cwd); err != nil {
			fmt.Println(err)
			return
		}
	},
}

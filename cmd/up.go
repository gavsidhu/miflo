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
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:     "up",
	Short:   "Apply migrations",
	Long:    "The up command applies all pending migrations in the migrations folder.",
	Args:    cobra.NoArgs,
	Example: "miflo up",
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

		if err := miflo.ApplyMigrations(database, ctx, cwd); err != nil {
			fmt.Println(err)
			return
		}

	},
}

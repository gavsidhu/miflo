package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/cli"
	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/miflo"
	"github.com/joho/godotenv"
)

var (
	showRevertHelp bool
)

func init() {
	revertCmd.Flags.Usage = printRevertHelp
	revertCmd.Flags.BoolVar(&showRevertHelp, "h", false, "Show help information for miflo revert")
	rootCmd.AddCommand(&revertCmd)
}

var revertCmd = cli.Command{
	Name:        "revert",
	Description: "Rollback most recent migration",
	Flags:       flag.NewFlagSet("revert", flag.ExitOnError),
	Run: func(cmd *cli.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}

		cmd.Flags.Parse(args)

		if showRevertHelp {
			cmd.Flags.Usage()
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

func printRevertHelp() {
	fmt.Println(`
Revert last migration.

Usage: miflo revert`)
}

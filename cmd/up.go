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
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	showUpHelp bool
)

func init() {
	upCmd.Flags.Usage = printUpHelp
	upCmd.Flags.BoolVar(&showUpHelp, "h", false, "Show help information for miflo up")
	rootCmd.AddCommand(&upCmd)
}

var upCmd = cli.Command{
	Name:        "up",
	Description: "Apply all pending migrations",
	Flags:       flag.NewFlagSet("apply", flag.ExitOnError),
	Run: func(cmd *cli.Command, args []string) {

		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}

		cmd.Flags.Parse(args)

		if showUpHelp {
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

		if err := miflo.ApplyMigrations(database, ctx, cwd); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func printUpHelp() {
	fmt.Println(`
Apply pending migrations.

Usage: miflo up`)
}

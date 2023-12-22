package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/cli"
	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/miflo"
	_ "github.com/lib/pq"
)

var (
	showListHelp bool
)

func init() {
	listCmd.Flags.Usage = printListHelp
	listCmd.Flags.BoolVar(&showListHelp, "h", false, "Show help information for miflo list")
	rootCmd.AddCommand(&listCmd)
}

var listCmd = cli.Command{
	Name:        "list",
	Description: "List all pending migrations",
	Flags:       flag.NewFlagSet("list", flag.ExitOnError),
	Run: func(cmd *cli.Command, args []string) {

		cmd.Flags.Parse(args)

		if showListHelp {
			cmd.Flags.Usage()
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

func printListHelp() {
	fmt.Println(`
List pending migrations.

Usage: miflo list`)
}

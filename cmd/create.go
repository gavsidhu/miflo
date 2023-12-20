package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/gavsidhu/miflo/internal/cli"
	"github.com/gavsidhu/miflo/internal/helpers"
	"github.com/gavsidhu/miflo/internal/miflo"
)

var (
	showCreateHelp bool
)

func init() {
	createCmd.Flags.BoolVar(&showCreateHelp, "h", false, "Show help information for miflo create")
	createCmd.Flags.Usage = printCreateHelp
	rootCmd.AddCommand(&createCmd)
}

var createCmd = cli.Command{
	Name:        "create",
	Description: "Create a new migration",
	Flags:       flag.NewFlagSet("create", flag.ExitOnError),
	Run: func(cmd *cli.Command, args []string) {

		cmd.Flags.Parse(args)

		if showCreateHelp {
			cmd.Flags.Usage()
			return
		}

		if len(args) < 1 {
			fmt.Printf("%sMigration name argument missing.%s\n", helpers.ColorYellow, helpers.ColorReset)
			cmd.Flags.Usage()
			return
		}

		cwd, err := os.Getwd()

		if err != nil {
			fmt.Println("error getting current working directory: ", err)
		}

		migrationName := args[0]

		if err := miflo.CreateMigration(migrationName, cwd); err != nil {
			fmt.Println("error creating migration: ", err)
		}

	},
}

func printCreateHelp() {
	fmt.Println(`
Create a new migration.

Usage: miflo create <migration_name>

Example: miflo create add_users_table`)
}

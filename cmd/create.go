package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

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

		var migrationsDirExists bool

		_, err = os.Stat(path.Join(cwd, "migrations"))

		if err != nil {
			if os.IsNotExist(err) {
				migrationsDirExists = false
			} else {
				fmt.Println("Error checking migrations directory:", err)
				return
			}
		} else {
			migrationsDirExists = true
		}

		if !migrationsDirExists {
			createDir := helpers.PromptForConfirmation("Migrations folder does not exist. Would you like to create it?")

			if createDir {
				os.Mkdir(path.Join(cwd, "migrations"), os.ModePerm)
			} else {
				return
			}
		}

		migrationName := args[0]
		timestamp := time.Now().Unix()

		if err := miflo.CreateMigration(migrationName, cwd, timestamp); err != nil {
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

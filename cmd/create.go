package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gavsidhu/miflo/internal/helpers"
	"github.com/gavsidhu/miflo/internal/miflo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create <migration name>",
	Short:   "Create a migration",
	Long:    "The create command creates a new migration file in the migrations folder. If there is no migration folder you will be prompted to create one in your root directory.",
	Args:    cobra.ExactArgs(1),
	Example: `miflo create setup_db_tables`,
	Run: func(cmd *cobra.Command, args []string) {

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

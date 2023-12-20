package miflo

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gavsidhu/miflo/internal/helpers"
)

func CreateMigration(migrationName string, cwd string) error {

	var migrationsDirExists bool

	_, err := os.Stat(path.Join(cwd, "migrations"))

	if err != nil {
		if os.IsNotExist(err) {
			migrationsDirExists = false
		} else {
			fmt.Println("Error checking migrations directory:", err)
			return err
		}
	} else {
		migrationsDirExists = true
	}

	if !migrationsDirExists {
		createDir := helpers.PromptForConfirmation("Migrations folder does not exist. Would you like to create it?")

		if createDir {
			os.Mkdir(path.Join(cwd, "migrations"), os.ModePerm)
		} else {
			return err
		}
	}

	timestamp := time.Now().Unix()

	pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, migrationName))

	os.Mkdir(pathName, os.ModePerm)
	os.Create(path.Join(pathName, "up.sql"))
	os.Create(path.Join(pathName, "down.sql"))

	fmt.Println("Migration created successfully")

	return nil

}

package miflo

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/gavsidhu/miflo/internal/helpers"
)

func CreateMigration(migrationName string, cwd string, timestamp int64) error {

	validMigrationName := helpers.IsValidMigrationName(migrationName)

	if !validMigrationName {
		return errors.New("invalid migration name")
	}

	pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, migrationName))

	if err := os.Mkdir(pathName, os.ModePerm); err != nil {
		return err
	}
	if _, err := os.Create(path.Join(pathName, "up.sql")); err != nil {
		return err
	}

	if _, err := os.Create(path.Join(pathName, "down.sql")); err != nil {
		return err
	}

	fmt.Println("Migration created successfully")

	return nil

}

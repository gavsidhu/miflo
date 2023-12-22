package miflo

import (
	"fmt"
	"os"
	"path"
)

func CreateMigration(migrationName string, cwd string, timestamp int64) error {

	pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, migrationName))

	os.Mkdir(pathName, os.ModePerm)
	os.Create(path.Join(pathName, "up.sql"))
	os.Create(path.Join(pathName, "down.sql"))

	fmt.Println("Migration created successfully")

	return nil

}

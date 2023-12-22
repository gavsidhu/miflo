package miflo

import (
	"fmt"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/helpers"
)

func ListPendingMigrations(db database.Database, cwd string) error {
	dirMigrations, err := helpers.GetDirMigrations(cwd)
	if err != nil {
		return err
	}

	appliedMigrationsRows, err := db.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("error getting applied migrations: %w", err)
	}

	defer appliedMigrationsRows.Close()

	appliedMigrations, err := helpers.GetAppliedMigrationNames(appliedMigrationsRows)
	if err != nil {
		return fmt.Errorf("error getting applied migration names: %w", err)
	}

	var pendingMigrations []string
	for _, migration := range dirMigrations {
		if !helpers.Contains(appliedMigrations, migration) {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	if len(pendingMigrations) < 1 {
		fmt.Println("No pending migrations")
		return nil
	}

	fmt.Println("Pending migrations:")

	for _, pending := range pendingMigrations {
		fmt.Println(helpers.ColorYellow, pending, helpers.ColorReset)
	}

	return nil
}

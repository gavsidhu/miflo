package miflo

import (
	"fmt"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/helpers"
)

func ListPendingMigrations(db database.Database) error {
	dirMigrations := helpers.GetDirMigrations()

	appliedMigrationsRows, err := db.GetAppliedMigrations()
	if err != nil {
		fmt.Println("error getting applied mitgrations:", err)
	}

	defer appliedMigrationsRows.Close()

	appliedMigrations, err := helpers.GetAppliedMigrationNames(appliedMigrationsRows)
	if err != nil {
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

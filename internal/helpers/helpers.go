package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ErrAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetAppliedMigrationNames(rows *sql.Rows) ([]string, error) {
	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations = append(migrations, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return migrations, nil
}

func GetDirMigrations(cwd string) ([]string, error) {
	entries, err := os.ReadDir(path.Join(cwd, "migrations"))
	if err != nil {
		return nil, err
	}
	var migrations []string

	for _, entry := range entries {
		if entry.IsDir() {
			migrations = append(migrations, entry.Name())
		}
	}

	return migrations, nil
}

func SortDirMigrations(migrations []string, ascending bool) {
	sort.Slice(migrations, func(i, j int) bool {
		timeI, _ := strconv.ParseInt(strings.Split(migrations[i], "_")[0], 10, 64)
		timeJ, _ := strconv.ParseInt(strings.Split(migrations[j], "_")[0], 10, 64)
		if ascending {
			return timeI < timeJ
		}
		return timeJ < timeI
	})
}

func IsValidMigrationName(migrationName string) bool {
	validNamePattern := regexp.MustCompile(`^[A-Za-z_]+$`)

	return validNamePattern.MatchString(migrationName)
}

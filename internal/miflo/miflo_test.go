package miflo_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/miflo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/stretchr/testify/assert"
)

type dbTestCase struct {
	name        string
	databaseURL string
}

func dbTestCases() []dbTestCase {
	return []dbTestCase{
		{
			name:        "SQLite",
			databaseURL: os.Getenv("SQLITE_TEST_DATABASE_URL"),
		},
		{
			name:        "PostgreSQL",
			databaseURL: os.Getenv("POSTGRES_TEST_DATABASE_URL"),
		},
		{
			name:        "libSQL",
			databaseURL: os.Getenv("LIBSQL_TEST_DATABASE_URL"),
		},
	}
}

func setupEnv(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	envPath := path.Join(cwd, "../../.env")

	err = godotenv.Load(envPath)

	if err != nil {
		t.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}
}

func newTestDatabase(t *testing.T, databaseURL string) database.Database {

	db, err := database.NewDatabase(databaseURL)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
		return nil
	}

	defer t.Cleanup(func() {
		setupEnv(t)
		err := os.Remove(strings.TrimPrefix(os.Getenv("SQLITE_TEST_DATABASE_URL"), "sqlite:"))
		if err != nil {
			t.Logf("Failed to delete test database file: %v", err)
		}
	})

	return db
}

func TestCreateMigrations(t *testing.T) {
	setupEnv(t)
	tests := []struct {
		name                      string
		migrationName             string
		expectMigrationFileCreate bool
		expectedError             bool
		setupFunc                 func(string)
		cleanupFunc               func(string)
	}{
		{
			name:                      "MigrationsDirExists",
			migrationName:             "initial_migration",
			expectMigrationFileCreate: true,
			expectedError:             false,
			setupFunc: func(cwd string) {
				os.Mkdir(path.Join(cwd, "migrations"), os.ModePerm)
			},
			cleanupFunc: func(cwd string) {
				migrationsDir := path.Join(cwd, "migrations")
				err := os.RemoveAll(migrationsDir)
				if err != nil {
					fmt.Printf("Failed to remove migrations directory: %v\n", err)
				}
			},
		},
		{
			name:                      "MigrationsDirNotExists",
			migrationName:             "initial_migration",
			expectMigrationFileCreate: false,
			expectedError:             true,
			cleanupFunc: func(cwd string) {
				migrationsDir := path.Join(cwd, "migrations")
				err := os.RemoveAll(migrationsDir)
				if err != nil {
					fmt.Printf("Failed to remove migrations directory: %v\n", err)
				}
			},
		},
		{
			name:                      "InvalidMigrationName",
			migrationName:             "invalid-name",
			expectMigrationFileCreate: false,
			expectedError:             true,
			setupFunc: func(cwd string) {
				os.Mkdir(path.Join(cwd, "migrations"), os.ModePerm)
			},
			cleanupFunc: func(cwd string) {
				migrationsDir := path.Join(cwd, "migrations")
				err := os.RemoveAll(migrationsDir)
				if err != nil {
					fmt.Printf("Failed to remove migrations directory: %v\n", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cwd := "../../test-migrations"

			if tt.setupFunc != nil {
				tt.setupFunc(cwd)
			}

			timestamp := time.Now().Unix()

			err := miflo.CreateMigration(tt.migrationName, cwd, timestamp)

			if tt.expectedError {
				assert.Error(t, err, "Expected an error for %s", tt.name)
			} else {
				assert.NoError(t, err, "Expected no error for %s", tt.name)
				if tt.expectMigrationFileCreate {
					timestamp := time.Now().Unix()
					dirName := fmt.Sprintf("%d_%s", timestamp, tt.migrationName)
					_, err = os.Stat(path.Join(cwd, "migrations", dirName))
					assert.NoError(t, err)
					_, err = os.Stat(path.Join(cwd, "migrations", dirName, "up.sql"))
					assert.NoError(t, err)
					_, err = os.Stat(path.Join(cwd, "migrations", dirName, "down.sql"))
					assert.NoError(t, err)
				}
			}

			if tt.cleanupFunc != nil {
				tt.cleanupFunc(cwd)
			}
		})
	}

}

func TestApplyMigration(t *testing.T) {
	setupEnv(t)

	tests := []struct {
		name                      string
		expectedError             bool
		setupFunc                 func(db database.Database, ctx context.Context, cwd string, dbName string) error
		cleanupFunc               func(db database.Database, ctx context.Context, cwd string) error
		postMigrationVerification func(t *testing.T, db database.Database, ctx context.Context, cwd string)
	}{
		{
			name:          "NoPendingigrations",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				if err := os.Mkdir(path.Join(cwd, "migrations"), os.ModePerm); err != nil {
					return err
				}
				return nil
			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:          "ApplyPendingMigrations",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", time.Now().Unix(), "test_migration"))
				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				upSqlPath := path.Join(pathName, "up.sql")
				file, err := os.OpenFile(upSqlPath, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %v", upSqlPath, err)
				}
				defer file.Close()

				sqlStatement := "CREATE TABLE IF NOT EXISTS miflo_test (id INT PRIMARY KEY, name TEXT);"
				if _, err := fmt.Fprintln(file, sqlStatement); err != nil {
					return fmt.Errorf("failed to write to file %s: %v", upSqlPath, err)
				}

				if _, err := os.Create(path.Join(pathName, "down.sql")); err != nil {
					return err
				}
				return nil

			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				query := "DROP TABLE IF EXISTS miflo_test"

				_, err := db.ExecContext(ctx, query)
				if err != nil {
					return err
				}

				query = "DELETE from migrations"

				_, err = db.ExecContext(ctx, query)
				if err != nil {
					return err
				}

				return nil
			},
			postMigrationVerification: func(t *testing.T, db database.Database, ctx context.Context, cwd string) {
				_, err := db.ExecContext(ctx, "SELECT 1 FROM miflo_test LIMIT 1")
				assert.NoError(t, err, "Migration should have created miflo_test")
			},
		},
		{
			name:          "MigrationAlreadyApplied",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				timestamp := time.Now().Unix()
				pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, "test_migration"))

				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				upSqlPath := path.Join(pathName, "up.sql")
				file, err := os.OpenFile(upSqlPath, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %v", upSqlPath, err)
				}
				defer file.Close()

				sqlStatement := "CREATE TABLE IF NOT EXISTS miflo_test (id INT PRIMARY KEY, name TEXT);"
				if _, err := fmt.Fprintln(file, sqlStatement); err != nil {
					return fmt.Errorf("failed to write to file %s: %v", upSqlPath, err)
				}

				if _, err := os.Create(path.Join(pathName, "down.sql")); err != nil {
					return err
				}
				if dbName == "PostgreSQL" {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES($1, $2, $3)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}
				} else {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES(?, ?, ?)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}

				}

				return nil

			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				query := "DROP TABLE IF EXISTS miflo_test"

				_, err := db.ExecContext(ctx, query)
				if err != nil {
					return err
				}

				query = "DELETE from migrations"

				_, err = db.ExecContext(ctx, query)
				if err != nil {
					return err
				}

				return nil
			},
			postMigrationVerification: func(t *testing.T, db database.Database, ctx context.Context, cwd string) {
				_, err := db.ExecContext(ctx, "SELECT 1 FROM miflo_test LIMIT 1")
				assert.Error(t, err, "DB should not have miflo_test table")

				rows, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM migrations")
				assert.NoError(t, err)
				defer rows.Close()

				var count int
				if rows.Next() {
					err := rows.Scan(&count)
					assert.NoError(t, err)
					assert.NotZero(t, count, "migrations should have rows")
				}

			},
		},
	}

	for _, dbCase := range dbTestCases() {

		db := newTestDatabase(t, dbCase.databaseURL)
		if db == nil {
			t.Fatal("error setting up test database")
		}

		ctx := context.Background()
		cwd := "../../test-migrations"

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.setupFunc(db, ctx, cwd, dbCase.name)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}

				err = miflo.ApplyMigrations(db, ctx, cwd)

				if tt.postMigrationVerification != nil {
					tt.postMigrationVerification(t, db, ctx, cwd)
				}

				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				err = tt.cleanupFunc(db, ctx, cwd)
				assert.NoError(t, err)
			})
		}
	}
}

func TestRevertMigration(t *testing.T) {
	setupEnv(t)

	tests := []struct {
		name                      string
		expectedError             bool
		setupFunc                 func(db database.Database, ctx context.Context, cwd string, dbName string) error
		cleanupFunc               func(db database.Database, ctx context.Context, cwd string) error
		postMigrationVerification func(t *testing.T, db database.Database, ctx context.Context, cwd string)
	}{
		{
			name:          "RevertAppliedMigrations",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {

				if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS miflo_test (id INT PRIMARY KEY, name TEXT);"); err != nil {
					return err
				}

				timestamp := time.Now().Unix()

				if dbName == "PostgreSQL" {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES($1, $2, $3)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}
				} else {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES(?, ?, ?)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}

				}

				pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, "test_migration"))
				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				downFile := path.Join(pathName, "down.sql")
				file, err := os.OpenFile(downFile, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return fmt.Errorf("failed to open file %s: %v", downFile, err)
				}
				defer file.Close()

				sqlStatement := "DROP TABLE IF EXISTS miflo_test"

				if _, err := fmt.Fprintln(file, sqlStatement); err != nil {
					return fmt.Errorf("failed to write to file %s: %v", downFile, err)
				}

				if _, err := os.Create(path.Join(pathName, "up.sql")); err != nil {
					return err
				}
				return nil

			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				return nil
			},
			postMigrationVerification: func(t *testing.T, db database.Database, ctx context.Context, cwd string) {
				_, err := db.ExecContext(ctx, "SELECT 1 FROM miflo_test LIMIT 1")
				assert.Error(t, err, "DB should not have miflo_test table")

				rows, err := db.QueryContext(ctx, "SELECT * FROM migrations LIMIT 1")
				assert.NoError(t, err, "Query should not return an error")
				defer rows.Close()

				assert.False(t, rows.Next(), "migrations table should be empty")

			},
		},
		{
			name:          "RevertUnappliedMigration",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				if dbName == "PostgreSQL" {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES($1, $2, $3)", fmt.Sprintf("%d_%s", time.Now().Unix(), "test_migration"), 1, false); err != nil {
						return err
					}
				} else {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES(?, ?, ?)", fmt.Sprintf("%d_%s", time.Now().Unix(), "test_migration"), 1, false); err != nil {
						return err
					}

				}
				return nil
			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				if _, err := db.ExecContext(ctx, "DELETE FROM migrations"); err != nil {
					return err
				}

				return nil
			},
			postMigrationVerification: func(t *testing.T, db database.Database, ctx context.Context, cwd string) {
				rows, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM migrations")
				assert.NoError(t, err)
				defer rows.Close()

				var count int
				if rows.Next() {
					err := rows.Scan(&count)
					assert.NoError(t, err)
					assert.NotZero(t, count, "migrations should have rows")
				}

			},
		},
		{
			name:          "NoMigrationsToRevert",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				return nil
			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				return nil
			},
		},
	}
	for _, dbCase := range dbTestCases() {

		db := newTestDatabase(t, dbCase.databaseURL)
		if db == nil {
			t.Fatal("error setting up test database")
		}

		ctx := context.Background()
		cwd := "../../test-migrations"

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.setupFunc(db, ctx, cwd, dbCase.name)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}

				err = miflo.RevertMigrations(db, ctx, cwd)

				if tt.postMigrationVerification != nil {
					tt.postMigrationVerification(t, db, ctx, cwd)
				}

				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				err = tt.cleanupFunc(db, ctx, cwd)
				assert.NoError(t, err)
			})
		}
	}
}

func TestListMigrations(t *testing.T) {
	setupEnv(t)

	tests := []struct {
		name          string
		expectedError bool
		setupFunc     func(db database.Database, ctx context.Context, cwd string, dbName string) error
		cleanupFunc   func(db database.Database, ctx context.Context, cwd string) error
	}{
		{
			name:          "ListPendingMigrations",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", time.Now().Unix(), "test_migration"))
				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				return nil
			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {
				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				return nil
			},
		},
		{
			name:          "AllMigrationsApplied",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd string, dbName string) error {
				timestamp := time.Now().Unix()

				if dbName == "PostgreSQL" {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES($1, $2, $3)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}
				} else {

					if _, err := db.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES(?, ?, ?)", fmt.Sprintf("%d_%s", timestamp, "test_migration"), 1, true); err != nil {
						return err
					}

				}

				pathName := path.Join(cwd, "migrations", fmt.Sprintf("%d_%s", timestamp, "test_migration"))
				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				return nil

			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {

				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				if _, err := db.ExecContext(ctx, "DELETE FROM migrations"); err != nil {
					return err
				}

				return nil
			},
		},
		{
			name:          "NoPendingMigrations",
			expectedError: false,
			setupFunc: func(db database.Database, ctx context.Context, cwd, dbName string) error {
				pathName := path.Join(cwd, "migrations")
				if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
					return err
				}

				return nil
			},
			cleanupFunc: func(db database.Database, ctx context.Context, cwd string) error {

				migrationsDir := path.Join(cwd, "migrations")
				if err := os.RemoveAll(migrationsDir); err != nil {
					return err
				}

				return nil
			},
		},
	}

	for _, dbCase := range dbTestCases() {

		db := newTestDatabase(t, dbCase.databaseURL)
		if db == nil {
			t.Fatal("error setting up test database")
		}

		ctx := context.Background()
		cwd := "../../test-migrations"

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setupFunc != nil {
					err := tt.setupFunc(db, ctx, cwd, dbCase.name)
					if err != nil {
						t.Fatalf("Setup failed: %v", err)
					}
				}

				err := miflo.ListPendingMigrations(db, cwd)

				if tt.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				if tt.cleanupFunc != nil {
					err = tt.cleanupFunc(db, ctx, cwd)
					assert.NoError(t, err)
				}

			})
		}
	}
}

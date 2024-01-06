package database

import (
	"os"
	"path"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func TestNewDatabase(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	envPath := path.Join(cwd, "../../.env")

	err = godotenv.Load(envPath)

	if err != nil {
		t.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	tests := []struct {
		name    string
		dbURL   string
		wantErr bool
	}{
		{"Valid SQLite", "sqlite:./test.db", false},
		{"Invalid SQLite", "./test.db", true},
		{"Invalid SQLite", "sqlite://./test.db", true},
		// Need to run docker container
		{"Valid PostgreSQL", "postgresql://testuser:testpassword@localhost:5432/testdb?sslmode=disable", false},
		{"Valid PostgreSQL", "postgres://testuser:testpassword@localhost:5432/testdb?sslmode=disable", false},
		{"Invalid PostgreSQL URL", "user:password@localhost/dbname", true},
		{"Valid libsql ", os.Getenv("LIBSQL_TEST_DATABASE_URL"), false},
		{"Valid libsql", "http://127.0.0.1:8080", false},
		{"Invalid libsql ", "db.turso.io", true},
		{"Invalid libsql", "127.0.0.1:8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatabase(tt.dbURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Cleanup(func() {
		err := os.Remove("./test.db")
		if err != nil {
			t.Logf("Failed to delete test database file: %v", err)
		}
	})
}

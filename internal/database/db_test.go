package database

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestNewDatabase(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDatabase(tt.dbURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

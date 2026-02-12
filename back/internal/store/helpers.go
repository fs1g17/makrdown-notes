package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	return db
}

func TruncateTables(t *testing.T, db *sql.DB) {
	tables := []string{"tokens", "notes", "folders", "users"} // order matters (FK constraints)
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
}

func CompareUsers(t *testing.T, expectedUser *User, actualUser *User) {
	t.Helper()
	assert.Equal(t, expectedUser.ID, actualUser.ID)
	assert.Equal(t, expectedUser.Username, actualUser.Username)
	assert.Equal(t, expectedUser.Email, actualUser.Email)
	assert.Equal(t, expectedUser.PasswordHash.hash, actualUser.PasswordHash.hash)
	assert.Equal(t, expectedUser.CreatedAt, actualUser.CreatedAt)
	assert.Equal(t, expectedUser.UpdatedAt, actualUser.UpdatedAt)
}

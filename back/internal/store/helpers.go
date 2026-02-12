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

func CompareFolders(t *testing.T, f1 *Folder, f2 *Folder) {
	t.Helper()
	assert.Equal(t, f1.ID, f2.ID)
	assert.Equal(t, f1.UserID, f2.UserID)
	assert.Equal(t, f1.ParentID, f2.ParentID)
	assert.Equal(t, f1.Name, f2.Name)
	assert.Equal(t, f1.CreatedAt, f2.CreatedAt)
	assert.Equal(t, f1.UpdatedAt, f2.UpdatedAt)
}

func CreateRootFolder(t *testing.T, db *sql.DB, folderStore PostgresFoldersStore, user *User) int64 {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	rootFolderId, err := folderStore.CreateFolderTx(tx, user.ID, nil, "root")
	if err != nil {
		t.Fatalf("failed to create root folder: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit tx: %v", err)
	}

	return rootFolderId
}

func CreateTestUser(t *testing.T, db *sql.DB, userStore UserStore, username string, email string, password string) (*User, error) {
	t.Helper()

	user := &User{
		Username: username,
		Email:    email,
	}

	if err := user.PasswordHash.Set(password); err != nil {
		t.Fatalf("failed to set password: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	err = userStore.CreateUser(tx, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	return user, nil
}

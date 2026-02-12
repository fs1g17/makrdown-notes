package service

import (
	"markdown-notes/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db := store.SetupTestDB(t)
	store.TruncateTables(t, db)
	userStore := store.NewPostgresUserStore(db)
	folderStore := store.NewPostgresFoldersStore(db)
	registerUserService := NewRegisterUserService(db, userStore, folderStore)

	user := &store.User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}

	if err := user.PasswordHash.Set("Password"); err != nil {
		t.Fatalf("failed to set password: %v", err)
	}

	rootFolderId, err := registerUserService.RegisterUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID, "ID should be populated by RETURNING clause")
	assert.NotZero(t, rootFolderId)

	dbUser, err := userStore.GetUserByUsername(user.Username)
	assert.NoError(t, err)
	store.CompareUsers(t, user, dbUser)

	dbRootFolderId, err := folderStore.GetRootFolder(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, rootFolderId, dbRootFolderId)
}

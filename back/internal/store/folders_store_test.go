package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createSubFolder(t *testing.T, db *sql.DB, folderStore PostgresFoldersStore, user *User, parent_id int64, name string) *Folder {
	t.Helper()

	folder, err := folderStore.CreateFolder(user.ID, parent_id, name)
	if err != nil {
		t.Fatalf("failed to create subfolder: %v", err)
	}

	return folder
}

func TestCreateFolder(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	folderStore := NewPostgresFoldersStore(db)
	userStore := NewPostgresUserStore(db)

	user, err := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	assert.NoError(t, err)

	var rootFolderId int64
	subfolderName := "subfolder"
	t.Run("successfully creates root folder", func(t *testing.T) {
		rootFolderId = CreateRootFolder(t, db, *folderStore, user)

		query := `
		SELECT id, user_id, parent_id, name, created_at, updated_at
		FROM folders 
		WHERE id = $1;
		`

		var rootFolder Folder
		err = db.QueryRow(query, rootFolderId).Scan(&rootFolder.ID, &rootFolder.UserID, &rootFolder.ParentID, &rootFolder.Name, &rootFolder.CreatedAt, &rootFolder.UpdatedAt)
		assert.NoError(t, err)

		assert.Equal(t, rootFolderId, rootFolder.ID)
		assert.Equal(t, user.ID, rootFolder.UserID)
		assert.Nil(t, rootFolder.ParentID)
		assert.Equal(t, "root", rootFolder.Name)
	})

	t.Run("successfully create subfolder", func(t *testing.T) {
		_, err := folderStore.CreateFolder(user.ID, rootFolderId, subfolderName)
		assert.NoError(t, err)

		query := `
		SELECT id, user_id, parent_id, name, created_at, updated_at
		FROM folders 
		WHERE name = $1;
		`

		var subfolder Folder
		err = db.QueryRow(query, subfolderName).Scan(&subfolder.ID, &subfolder.UserID, &subfolder.ParentID, &subfolder.Name, &subfolder.CreatedAt, &subfolder.UpdatedAt)
		assert.NoError(t, err)

		assert.Equal(t, user.ID, subfolder.UserID)
		assert.Equal(t, rootFolderId, *subfolder.ParentID)
		assert.Equal(t, subfolderName, subfolder.Name)
	})

	t.Run("fail to create subfolder with clashing name", func(t *testing.T) {
		subfolder, err := folderStore.CreateFolder(user.ID, rootFolderId, subfolderName)
		assert.Nil(t, subfolder)
		assert.Error(t, err)
	})
}

func TestGetRootFolder(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	folderStore := NewPostgresFoldersStore(db)
	userStore := NewPostgresUserStore(db)

	user, err := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	assert.NoError(t, err)

	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	t.Run("gets root folder for existing user", func(t *testing.T) {
		dbFolderId, err := folderStore.GetRootFolder(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, rootFolderId, dbFolderId)
	})

	t.Run("fails to get root folder for non-existent user", func(t *testing.T) {
		dbFolderId, err := folderStore.GetRootFolder(72)
		assert.Error(t, err)
		assert.Equal(t, int64(0), dbFolderId)
	})
}

func TestUserOwnsFolder(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	folderStore := NewPostgresFoldersStore(db)
	userStore := NewPostgresUserStore(db)

	user, err := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	assert.NoError(t, err)
	user2, err := CreateTestUser(t, db, userStore, "Theo2", "example@gmail.com", "Password")
	assert.NoError(t, err)

	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	t.Run("returns true when user owns folder", func(t *testing.T) {
		owns, err := folderStore.UserOwnsFolder(user.ID, rootFolderId)
		assert.NoError(t, err)
		assert.True(t, owns)
	})

	t.Run("returns false when user doesn't own folder", func(t *testing.T) {
		owns, err := folderStore.UserOwnsFolder(user2.ID, rootFolderId)
		assert.Error(t, err)
		assert.False(t, owns)
	})
}

func TestGetSubFolders(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	folderStore := NewPostgresFoldersStore(db)
	userStore := NewPostgresUserStore(db)

	user, err := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	assert.NoError(t, err)

	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	t.Run("returns empty array when no subfolders exist", func(t *testing.T) {
		folders, err := folderStore.GetSubFolders(user.ID, rootFolderId)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(folders))
	})

	t.Run("returns empty array when called with non-existent parent-id", func(t *testing.T) {
		folders, err := folderStore.GetSubFolders(user.ID, 72)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(folders))
	})

	t.Run("returns subfolders", func(t *testing.T) {
		subFolder := createSubFolder(t, db, *folderStore, user, rootFolderId, "subfolder")
		folders, err := folderStore.GetSubFolders(user.ID, rootFolderId)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(folders))

		CompareFolders(t, subFolder, &(folders[0]))
	})

	t.Run("does not return other user's subfolders", func(t *testing.T) {
		user2, _ := CreateTestUser(t, db, userStore, "Theo2", "example@gmail.com", "Password")
		folders, err := folderStore.GetSubFolders(user2.ID, rootFolderId)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(folders))
	})
}

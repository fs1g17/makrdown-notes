package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFolder(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	folderStore := NewPostgresFoldersStore(db)
	userStore := NewPostgresUserStore(db)

	user, err := createTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	assert.NoError(t, err)

	var rootFolderId int64
	subfolderName := "subfolder"
	t.Run("successfully creates root folder", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("failed to begin tx: %v", err)
		}

		rootFolderId, err = folderStore.CreateFolderTx(tx, user.ID, nil, "root")
		assert.NoError(t, err)
		tx.Commit()

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

package service

import (
	"markdown-notes/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFolderContent(t *testing.T) {
	db := store.SetupTestDB(t)
	store.TruncateTables(t, db)
	userStore := store.NewPostgresUserStore(db)
	notesStore := store.NewPostgresNotesStore(db)
	folderStore := store.NewPostgresFoldersStore(db)
	registerUserService := NewRegisterUserService(db, userStore, folderStore)
	folderContentsService := NewFolderContentsService(db, userStore, folderStore, notesStore)

	user := &store.User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}
	user.PasswordHash.Set("Password")

	user2 := &store.User{
		Username: "Other",
		Email:    "other@gmail.com",
	}
	user2.PasswordHash.Set("Password")

	rootFolderId, err := registerUserService.RegisterUser(user)
	assert.NoError(t, err)
	_, err = registerUserService.RegisterUser(user2)
	assert.NoError(t, err)

	t.Run("empty root folder content", func(t *testing.T) {
		folderContent, err := folderContentsService.GetFolderContent(user, rootFolderId)
		assert.NoError(t, err)

		assert.Equal(t, rootFolderId, folderContent.FolderID)
		assert.Equal(t, 0, len(folderContent.Folders))
		assert.Equal(t, 0, len(folderContent.Notes))
	})

	t.Run("returns folder content with subfolders and notes", func(t *testing.T) {
		// Create subfolder and note
		subfolder, err := folderContentsService.CreateSubFolder(user, rootFolderId, "subfolder")
		assert.NoError(t, err)
		note, err := folderContentsService.CreateNote(user, rootFolderId, "note title", "note content")
		assert.NoError(t, err)

		folderContent, err := folderContentsService.GetFolderContent(user, rootFolderId)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(folderContent.Folders))
		assert.Equal(t, subfolder.ID, folderContent.Folders[0].ID)
		assert.Equal(t, 1, len(folderContent.Notes))
		assert.Equal(t, note.ID, folderContent.Notes[0].ID)
	})

	t.Run("fails when user doesn't own folder", func(t *testing.T) {
		folderContent, err := folderContentsService.GetFolderContent(user2, rootFolderId)
		assert.Error(t, err)
		assert.Nil(t, folderContent)
	})
}

func TestCreateSubFolder(t *testing.T) {
	db := store.SetupTestDB(t)
	store.TruncateTables(t, db)
	userStore := store.NewPostgresUserStore(db)
	notesStore := store.NewPostgresNotesStore(db)
	folderStore := store.NewPostgresFoldersStore(db)
	registerUserService := NewRegisterUserService(db, userStore, folderStore)
	folderContentsService := NewFolderContentsService(db, userStore, folderStore, notesStore)

	user := &store.User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}
	user.PasswordHash.Set("Password")

	user2 := &store.User{
		Username: "Other",
		Email:    "other@gmail.com",
	}
	user2.PasswordHash.Set("Password")

	rootFolderId, err := registerUserService.RegisterUser(user)
	assert.NoError(t, err)
	_, err = registerUserService.RegisterUser(user2)
	assert.NoError(t, err)

	t.Run("creates subfolder in owned folder", func(t *testing.T) {
		subfolder, err := folderContentsService.CreateSubFolder(user, rootFolderId, "my-subfolder")
		assert.NoError(t, err)
		assert.NotNil(t, subfolder)
		assert.Equal(t, "my-subfolder", subfolder.Name)
		assert.Equal(t, rootFolderId, *subfolder.ParentID)
		assert.Equal(t, user.ID, subfolder.UserID)
	})

	t.Run("fails when user doesn't own parent folder", func(t *testing.T) {
		subfolder, err := folderContentsService.CreateSubFolder(user2, rootFolderId, "hacked-folder")
		assert.Error(t, err)
		assert.Nil(t, subfolder)
	})

	t.Run("fails for non-existent parent folder", func(t *testing.T) {
		subfolder, err := folderContentsService.CreateSubFolder(user, 9999, "orphan-folder")
		assert.Error(t, err)
		assert.Nil(t, subfolder)
	})
}

func TestCreateNote(t *testing.T) {
	db := store.SetupTestDB(t)
	store.TruncateTables(t, db)
	userStore := store.NewPostgresUserStore(db)
	notesStore := store.NewPostgresNotesStore(db)
	folderStore := store.NewPostgresFoldersStore(db)
	registerUserService := NewRegisterUserService(db, userStore, folderStore)
	folderContentsService := NewFolderContentsService(db, userStore, folderStore, notesStore)

	user := &store.User{
		Username: "Theo",
		Email:    "drumandbassbob@gmail.com",
	}
	user.PasswordHash.Set("Password")

	user2 := &store.User{
		Username: "Other",
		Email:    "other@gmail.com",
	}
	user2.PasswordHash.Set("Password")

	rootFolderId, err := registerUserService.RegisterUser(user)
	assert.NoError(t, err)
	_, err = registerUserService.RegisterUser(user2)
	assert.NoError(t, err)

	t.Run("creates note in specified folder", func(t *testing.T) {
		note, err := folderContentsService.CreateNote(user, rootFolderId, "my note", "my content")
		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, "my note", note.Title)
		assert.Equal(t, "my content", note.Note)
		assert.Equal(t, rootFolderId, note.FolderID)
	})

	t.Run("creates note in root folder when folder_id is 0", func(t *testing.T) {
		note, err := folderContentsService.CreateNote(user, 0, "root note", "root content")
		assert.NoError(t, err)
		assert.NotNil(t, note)
		assert.Equal(t, rootFolderId, note.FolderID)
	})

	t.Run("fails when user doesn't own folder", func(t *testing.T) {
		note, err := folderContentsService.CreateNote(user2, rootFolderId, "hacked note", "hacked content")
		assert.Error(t, err)
		assert.Nil(t, note)
	})

	t.Run("fails for non-existent folder", func(t *testing.T) {
		note, err := folderContentsService.CreateNote(user, 9999, "orphan note", "content")
		assert.Error(t, err)
		assert.Nil(t, note)
	})
}

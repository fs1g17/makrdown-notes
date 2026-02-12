package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)
	notesStore := NewPostgresNotesStore(db)
	folderStore := NewPostgresFoldersStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	t.Run("create note for existing user and folder", func(t *testing.T) {
		note, err := notesStore.CreateNote(user.ID, rootFolderId, "title", "note")
		assert.NoError(t, err)

		assert.NotZero(t, note.ID)
		assert.Equal(t, rootFolderId, note.FolderID)
		assert.Equal(t, "title", note.Title)
		assert.Equal(t, "note", note.Note)
	})

	t.Run("fails to create note with same name", func(t *testing.T) {
		note, err := notesStore.CreateNote(user.ID, rootFolderId, "title", "note")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateNote)
		assert.Nil(t, note)
	})

	t.Run("fails to create note for non-existent user id", func(t *testing.T) {
		note, err := notesStore.CreateNote(72, rootFolderId, "title", "note")
		assert.Error(t, err)
		assert.Nil(t, note)
	})

	t.Run("fails to create note for non-existent folder id", func(t *testing.T) {
		note, err := notesStore.CreateNote(user.ID, 72, "title", "note")
		assert.Error(t, err)
		assert.Nil(t, note)
	})
}

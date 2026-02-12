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

func TestGetNotesInFolder(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)
	notesStore := NewPostgresNotesStore(db)
	folderStore := NewPostgresFoldersStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	user2 := CreateTestUser(t, db, userStore, "Other", "other@gmail.com", "Password")
	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	t.Run("returns empty slice when no notes", func(t *testing.T) {
		notes, err := notesStore.GetNotesInFolder(user.ID, rootFolderId)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(notes))
	})

	t.Run("returns notes in folder", func(t *testing.T) {
		note1, err := notesStore.CreateNote(user.ID, rootFolderId, "note1", "content1")
		assert.NoError(t, err)
		note2, err := notesStore.CreateNote(user.ID, rootFolderId, "note2", "content2")
		assert.NoError(t, err)

		notes, err := notesStore.GetNotesInFolder(user.ID, rootFolderId)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(notes))
		assert.Equal(t, note1.ID, notes[0].ID)
		assert.Equal(t, note2.ID, notes[1].ID)
	})

	t.Run("does not return other user's notes", func(t *testing.T) {
		notes, err := notesStore.GetNotesInFolder(user2.ID, rootFolderId)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(notes))
	})
}

func TestGetNote(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)
	notesStore := NewPostgresNotesStore(db)
	folderStore := NewPostgresFoldersStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	user2 := CreateTestUser(t, db, userStore, "Other", "other@gmail.com", "Password")
	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	note, err := notesStore.CreateNote(user.ID, rootFolderId, "title", "content")
	assert.NoError(t, err)

	t.Run("returns note for valid user and note id", func(t *testing.T) {
		dbNote, err := notesStore.GetNote(user.ID, note.ID)
		assert.NoError(t, err)
		assert.Equal(t, note.ID, dbNote.ID)
		assert.Equal(t, note.Title, dbNote.Title)
		assert.Equal(t, note.Note, dbNote.Note)
		assert.Equal(t, note.FolderID, dbNote.FolderID)
	})

	t.Run("returns error for wrong user id", func(t *testing.T) {
		dbNote, err := notesStore.GetNote(user2.ID, note.ID)
		assert.Error(t, err)
		assert.Nil(t, dbNote)
	})

	t.Run("returns error for non-existent note", func(t *testing.T) {
		dbNote, err := notesStore.GetNote(user.ID, 9999)
		assert.Error(t, err)
		assert.Nil(t, dbNote)
	})
}

func TestUpdateNote(t *testing.T) {
	db := SetupTestDB(t)
	TruncateTables(t, db)
	userStore := NewPostgresUserStore(db)
	notesStore := NewPostgresNotesStore(db)
	folderStore := NewPostgresFoldersStore(db)

	user := CreateTestUser(t, db, userStore, "Theo", "drumandbassbob@gmail.com", "Password")
	user2 := CreateTestUser(t, db, userStore, "Other", "other@gmail.com", "Password")
	rootFolderId := CreateRootFolder(t, db, *folderStore, user)

	note, err := notesStore.CreateNote(user.ID, rootFolderId, "title", "original content")
	assert.NoError(t, err)

	t.Run("updates note content", func(t *testing.T) {
		updatedNote, err := notesStore.UpdateNote(user.ID, note.ID, "updated content")
		assert.NoError(t, err)
		assert.Equal(t, note.ID, updatedNote.ID)
		assert.Equal(t, note.Title, updatedNote.Title)
		assert.Equal(t, "updated content", updatedNote.Note)
	})

	t.Run("updates updated_at timestamp", func(t *testing.T) {
		originalNote, _ := notesStore.GetNote(user.ID, note.ID)
		updatedNote, err := notesStore.UpdateNote(user.ID, note.ID, "newer content")
		assert.NoError(t, err)
		assert.True(t, updatedNote.UpdatedAt.After(originalNote.CreatedAt) || updatedNote.UpdatedAt.Equal(originalNote.CreatedAt))
	})

	t.Run("fails for wrong user id", func(t *testing.T) {
		updatedNote, err := notesStore.UpdateNote(user2.ID, note.ID, "hacked content")
		assert.Error(t, err)
		assert.Nil(t, updatedNote)
	})

	t.Run("fails for non-existent note", func(t *testing.T) {
		updatedNote, err := notesStore.UpdateNote(user.ID, 9999, "content")
		assert.Error(t, err)
		assert.Nil(t, updatedNote)
	})
}

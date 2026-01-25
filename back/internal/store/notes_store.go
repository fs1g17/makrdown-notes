package store

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64     `json:"id"`
	FolderID  int64     `json:"folder_id"`
	Title     string    `json:"title"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresNotesStore struct {
	db *sql.DB
}

func NewPostgresNotesStore(db *sql.DB) *PostgresNotesStore {
	return &PostgresNotesStore{db: db}
}

type NotesStore interface {
	CreateNote(user_id int64, folder_id int64, title string, note string) (*Note, error)
	GetNotesInFolder(user_id int64, folder_id int64) ([]Note, error)
	GetNote(user_id int64, note_id int64) (*Note, error)
}

func (n *PostgresNotesStore) CreateNote(user_id int64, folder_id int64, title string, note string) (*Note, error) {
	query := `
	INSERT INTO notes (user_id, folder_id, title, note)
	VALUES ($1, $2, $3, $4)
	RETURNING id, folder_id, title, note, created_at, updated_at;
	`

	var dbNote Note
	err := n.db.QueryRow(query, user_id, folder_id, title, note).Scan(
		&dbNote.ID,
		&dbNote.FolderID,
		&dbNote.Title,
		&dbNote.Note,
		&dbNote.CreatedAt,
		&dbNote.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &dbNote, nil
}

func (n *PostgresNotesStore) GetNotesInFolder(user_id int64, folder_id int64) ([]Note, error) {
	query := `
	SELECT id, folder_id, title, note, created_at, updated_at
	FROM notes
	WHERE folder_id = $1
	ORDER BY updated_at;
	`

	rows, err := n.db.Query(query, folder_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		var note Note
		err = rows.Scan(
			&note.ID,
			&note.FolderID,
			&note.Title,
			&note.Note,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (n *PostgresNotesStore) GetNote(user_id int64, note_id int64) (*Note, error) {
	query := `
	SELECT id, folder_id, title, note, created_at, updated_at 
	FROM notes 
	WHERE user_id = $1 AND id = $2;
	`

	var dbNote Note
	err := n.db.QueryRow(query, user_id, note_id).Scan(
		&dbNote.ID,
		&dbNote.FolderID,
		&dbNote.Title,
		&dbNote.Note,
		&dbNote.CreatedAt,
		&dbNote.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &dbNote, nil
}

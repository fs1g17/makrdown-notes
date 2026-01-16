package store

import (
	"database/sql"
	"time"
)

type Note struct {
	ID        int64     `json:"id"`
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
	CreateNote(user_id int64, folder_id int64, title string, note string) (int64, error)
}

func (n *PostgresNotesStore) CreateNote(user_id int64, folder_id int64, title string, note string) (int64, error) {
	query := `
	INSERT INTO notes (user_id, folder_id, title, note)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	var note_id int64
	err := n.db.QueryRow(query, user_id, folder_id, title, note).Scan(&note_id)
	if err != nil {
		return 0, err
	}

	return note_id, nil
}

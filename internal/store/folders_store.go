package store

import (
	"database/sql"
	"time"
)

type Folder struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ParentID  *int64    `json:"parent_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresFoldersStore struct {
	db *sql.DB
}

func NewPostgresFoldersStore(db *sql.DB) *PostgresFoldersStore {
	return &PostgresFoldersStore{db: db}
}

type FoldersStore interface {
	CreateFolder(user_id int64, parent_id int64, name string) (Folder, error)
	CreateFolderTx(tx *sql.Tx, user_id int64, parent_id *int64, name string) (int64, error)
	GetRootFolder(user_id int64) (int64, error)
	UserOwnsFolder(user_id int64, folder_id int64) (bool, error)
	GetSubFolders(user_id int64, folder_id int64) ([]Folder, error)
}

func (f *PostgresFoldersStore) CreateFolder(user_id int64, parent_id int64, name string) (Folder, error) {
	query := `
	INSERT INTO folders (user_id, parent_id, name)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, parent_id, name, created_at, updated_at;
	`

	var folder Folder
	err := f.db.QueryRow(query, user_id, parent_id, name).Scan(
		&folder.ID,
		&folder.UserID,
		&folder.ParentID,
		&folder.Name,
		&folder.CreatedAt,
		&folder.UpdatedAt)
	if err != nil {
		return Folder{}, err
	}

	return folder, nil
}

func (f *PostgresFoldersStore) CreateFolderTx(tx *sql.Tx, user_id int64, parent_id *int64, name string) (int64, error) {
	query := `
	INSERT INTO folders (user_id, parent_id, name)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	var folder_id int64
	err := tx.QueryRow(query, user_id, parent_id, name).Scan(&folder_id)
	if err != nil {
		return 0, err
	}

	return folder_id, nil
}

func (f *PostgresFoldersStore) GetRootFolder(user_id int64) (int64, error) {
	query := `
	SELECT id 
	FROM folders 
	WHERE user_id = $1;
	`

	var folder_id int64
	err := f.db.QueryRow(query, user_id).Scan(&folder_id)
	if err != nil {
		return 0, err
	}

	return folder_id, nil
}

func (f *PostgresFoldersStore) UserOwnsFolder(user_id int64, folder_id int64) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1 
		FROM folders 
		WHERE user_id = $1 AND id = $2
	);
	`

	var exists bool
	err := f.db.QueryRow(query, user_id, folder_id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (f *PostgresFoldersStore) GetSubFolders(user_id int64, folder_id int64) ([]Folder, error) {
	query := `
	SELECT id, user_id, parent_id, name, created_at, updated_at
	FROM folders 
	WHERE user_id = $1 and parent_id = $2;
	`

	rows, err := f.db.Query(query, user_id, folder_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := []Folder{}

	for rows.Next() {
		var folder Folder
		err = rows.Scan(
			&folder.ID,
			&folder.UserID,
			&folder.ParentID,
			&folder.Name,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}

	return folders, nil
}

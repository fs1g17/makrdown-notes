package service

import (
	"database/sql"
	"markdown-notes/internal/store"
)

type RegisterUserService struct {
	db          *sql.DB
	userStore   store.UserStore
	folderStore store.FoldersStore
}

func NewRegisterUserService(db *sql.DB, userStore store.UserStore, folderStore store.FoldersStore) *RegisterUserService {
	return &RegisterUserService{
		db:          db,
		userStore:   userStore,
		folderStore: folderStore,
	}
}

type RegisterUserServiceI interface {
	RegisterUser(user *store.User) (int64, error)
}

func (s *RegisterUserService) RegisterUser(user *store.User) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = s.userStore.CreateUser(tx, user)
	if err != nil {
		return 0, err
	}

	folder_id, err := s.folderStore.CreateFolderTx(tx, user.ID, nil, "root")
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return folder_id, err
}

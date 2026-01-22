package service

import (
	"database/sql"
	"errors"
	"markdown-notes/internal/store"
)

type FolderContentsService struct {
	db          *sql.DB
	userStore   store.UserStore
	folderStore store.FoldersStore
	noteStore   store.NotesStore
}

func NewFolderContentsService(
	db *sql.DB,
	userStore store.UserStore,
	folderStore store.FoldersStore,
	noteStore store.NotesStore,
) *FolderContentsService {
	return &FolderContentsService{
		db,
		userStore,
		folderStore,
		noteStore,
	}
}

type FolderContent struct {
	Notes   []store.Note   `json:"notes"`
	Folders []store.Folder `json:"folders"`
}

type FolderContentsServiceI interface {
	GetFolderContent(user *store.User, folder_id int64) (FolderContent, error)
	CreateSubFolder(user *store.User, parent_id int64, name string) (store.Folder, error)
	CreateNote(user *store.User, folder_id int64, title string, note string) (store.Note, error)
}

func (f *FolderContentsService) GetFolderContent(user *store.User, folder_id int64) (FolderContent, error) {
	owns, err := f.folderStore.UserOwnsFolder(user.ID, folder_id)
	if err != nil {
		return FolderContent{}, err
	}

	if owns == false {
		return FolderContent{}, errors.New("unauthorized")
	}

	folders, err := f.folderStore.GetSubFolders(user.ID, folder_id)
	if err != nil {
		return FolderContent{}, err
	}

	notes, err := f.noteStore.GetNotesInFolder(user.ID, folder_id)
	if err != nil {
		return FolderContent{}, err
	}

	return FolderContent{
		Notes:   notes,
		Folders: folders,
	}, nil
}

func (f *FolderContentsService) CreateSubFolder(user *store.User, parent_id int64, name string) (store.Folder, error) {
	owns, err := f.folderStore.UserOwnsFolder(user.ID, parent_id)
	if err != nil {
		return store.Folder{}, err
	}

	if owns == false {
		return store.Folder{}, errors.New("unauthorized")
	}

	folder, err := f.folderStore.CreateFolder(user.ID, parent_id, name)
	if err != nil {
		return store.Folder{}, err
	}

	return folder, nil
}

func (f *FolderContentsService) CreateNote(user *store.User, folder_id int64, title string, note string) (store.Note, error) {
	owns, err := f.folderStore.UserOwnsFolder(user.ID, folder_id)
	if err != nil {
		return store.Note{}, err
	}

	if owns == false {
		return store.Note{}, errors.New("unauthorized")
	}

	dbNote, err := f.noteStore.CreateNote(user.ID, folder_id, title, note)
	if err != nil {
		return store.Note{}, err
	}

	return dbNote, nil
}

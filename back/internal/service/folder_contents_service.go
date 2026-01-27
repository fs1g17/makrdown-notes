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
	FolderID int64          `json:"folder_id"`
	Notes    []store.Note   `json:"notes"`
	Folders  []store.Folder `json:"folders"`
}

type FolderContentsServiceI interface {
	GetFolderContent(user *store.User, folder_id int64) (*FolderContent, error)
	CreateSubFolder(user *store.User, parent_id int64, name string) (*store.Folder, error)
	CreateNote(user *store.User, folder_id int64, title string, note string) (*store.Note, error)
}

func (f *FolderContentsService) GetFolderContent(user *store.User, folder_id int64) (*FolderContent, error) {
	owns, err := f.folderStore.UserOwnsFolder(user.ID, folder_id)
	if err != nil {
		return nil, err
	}

	if owns == false {
		return nil, errors.New("unauthorized")
	}

	folders, err := f.folderStore.GetSubFolders(user.ID, folder_id)
	if err != nil {
		return nil, err
	}

	notes, err := f.noteStore.GetNotesInFolder(user.ID, folder_id)
	if err != nil {
		return nil, err
	}

	return &FolderContent{
		FolderID: folder_id,
		Notes:    notes,
		Folders:  folders,
	}, nil
}

func (f *FolderContentsService) CreateSubFolder(user *store.User, parent_id int64, name string) (*store.Folder, error) {
	owns, err := f.folderStore.UserOwnsFolder(user.ID, parent_id)
	if err != nil {
		return nil, err
	}

	if owns == false {
		return nil, errors.New("unauthorized")
	}

	folder, err := f.folderStore.CreateFolder(user.ID, parent_id, name)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (f *FolderContentsService) CreateNote(user *store.User, folder_id int64, title string, note string) (*store.Note, error) {
	var use_folder_id int64
	if folder_id == 0 {
		// if folder_id is 0, get root folder id
		root_folder_id, err := f.folderStore.GetRootFolder(user.ID)
		if err != nil {
			return nil, err
		}
		use_folder_id = root_folder_id
	} else {
		// otherwise check if folder belongs to the user
		use_folder_id = folder_id
		owns, err := f.folderStore.UserOwnsFolder(user.ID, folder_id)
		if err != nil {
			return nil, err
		}

		if owns == false {
			return nil, errors.New("unauthorized")
		}
	}

	dbNote, err := f.noteStore.CreateNote(user.ID, use_folder_id, title, note)
	if err != nil {
		return nil, err
	}

	return dbNote, nil
}

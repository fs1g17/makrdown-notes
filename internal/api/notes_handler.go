package api

import (
	"errors"
	"log"
	"markdown-notes/internal/store"
	"markdown-notes/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotesHandler struct {
	notesStore store.NotesStore
	logger     *log.Logger
}

func NewNotesHandler(notesStore store.NotesStore, logger *log.Logger) *NotesHandler {
	return &NotesHandler{
		notesStore: notesStore,
		logger:     logger,
	}
}

type createNoteRequest struct {
	Title    string `json:"title"`
	Note     string `json:"note"`
	FolderID int64  `json:"folder_id"`
}

func (r *createNoteRequest) validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}

	if r.Note == "" {
		return errors.New("note is required")
	}

	if r.FolderID == 0 {
		return errors.New("folder_id is required")
	}

	return nil
}

func (h *NotesHandler) HandleCreateNote(c echo.Context) error {
	var req createNoteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := req.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	note_id, err := h.notesStore.CreateNote(int64(user.ID), req.FolderID, req.Title, req.Note)
	if err != nil {
		h.logger.Printf("Error: creating note")
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, utils.Envelope{"note_id": note_id})
}

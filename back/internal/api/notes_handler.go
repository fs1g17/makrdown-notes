package api

import (
	"database/sql"
	"errors"
	"log"
	"markdown-notes/internal/service"
	"markdown-notes/internal/store"
	"markdown-notes/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func httpStatusFromNoteError(err error) int {
	switch {
	case errors.Is(err, store.ErrDuplicateNote):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type NotesHandler struct {
	notesStore            store.NotesStore
	folderContentsService service.FolderContentsServiceI
	logger                *log.Logger
}

func NewNotesHandler(
	notesStore store.NotesStore,
	folderContentsService service.FolderContentsServiceI,
	logger *log.Logger,
) *NotesHandler {
	return &NotesHandler{
		notesStore:            notesStore,
		folderContentsService: folderContentsService,
		logger:                logger,
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
	note, err := h.folderContentsService.CreateNote(user, req.FolderID, req.Title, req.Note)
	if err != nil {
		h.logger.Printf("Error creating note: %v", err)
		return c.JSON(httpStatusFromNoteError(err), utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, note)
}

type getNoteRequest struct {
	NoteID int64 `param:"note_id"`
}

func (r *getNoteRequest) validate() error {
	if r.NoteID == 0 {
		return errors.New("note_id is required")
	}

	return nil
}

func (h *NotesHandler) HandleGetNote(c echo.Context) error {
	var req getNoteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := req.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	note, err := h.notesStore.GetNote(user.ID, req.NoteID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Printf("ERROR: user trying to access note %v", err)
			return c.JSON(http.StatusUnauthorized, utils.Envelope{"error": "note doesn't exist or you don't have access to it"})
		}
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, note)
}

type patchNoteRequest struct {
	NoteID int64  `param:"note_id"`
	Note   string `json:"note"`
}

func (r *patchNoteRequest) validate() error {
	if r.NoteID == 0 {
		return errors.New("note_id is required")
	}

	return nil
}

func (h *NotesHandler) HandlePatchNote(c echo.Context) error {
	var req patchNoteRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := req.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	note, err := h.notesStore.UpdateNote(user.ID, req.NoteID, req.Note)
	if err != nil {
		h.logger.Printf("ERROR: couldn't update the note %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, utils.Envelope{"note": note})
}

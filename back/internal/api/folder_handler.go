package api

import (
	"errors"
	"log"
	"markdown-notes/internal/service"
	"markdown-notes/internal/store"
	"markdown-notes/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func httpStatusFromFolderError(err error) int {
	switch {
	case errors.Is(err, store.ErrDuplicateFolder):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type FolderHandler struct {
	folderContentsService service.FolderContentsServiceI
	folderStore           store.FoldersStore
	logger                *log.Logger
}

func NewFolderHandler(
	folderContentsService service.FolderContentsServiceI,
	folderStore store.FoldersStore,
	logger *log.Logger,
) *FolderHandler {
	return &FolderHandler{
		folderContentsService: folderContentsService,
		folderStore:           folderStore,
		logger:                logger,
	}
}

type createFolderRequest struct {
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name"`
}

func (r *createFolderRequest) validate() error {
	if r.ParentID == 0 {
		return errors.New("parent_id is required")
	}

	if r.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (h *FolderHandler) HandleCreateFolder(c echo.Context) error {
	var req createFolderRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := req.validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	folder, err := h.folderContentsService.CreateSubFolder(user, req.ParentID, req.Name)
	if err != nil {
		h.logger.Printf("Error creating folder: %v", err)
		return c.JSON(httpStatusFromFolderError(err), utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, folder)
}

func (h *FolderHandler) GetRootFolderContent(c echo.Context) error {
	user := c.Get("user").(*store.User)
	root_folder_id, err := h.folderStore.GetRootFolder(user.ID)
	if err != nil {
		h.logger.Printf("Error: getting root folder id %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	folderContents, err := h.folderContentsService.GetFolderContent(user, root_folder_id)
	if err != nil {
		h.logger.Printf("Error: getting folder content %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, folderContents)
}

type getFolderContentRequest struct {
	FolderID int64 `param:"folder_id"`
}

func (r *getFolderContentRequest) validate() error {
	if r.FolderID == 0 {
		return errors.New("folder_id is required")
	}

	return nil
}

func (h *FolderHandler) GetFolderContent(c echo.Context) error {
	var req getFolderContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := req.validate()
	if err != nil {
		h.logger.Printf("Error: invalid request %v", err)
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	folderContents, err := h.folderContentsService.GetFolderContent(user, req.FolderID)
	if err != nil {
		h.logger.Printf("Error: getting folder content %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, folderContents)
}

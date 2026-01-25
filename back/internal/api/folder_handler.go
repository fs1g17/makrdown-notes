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
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, folder)
}

type getFolderContentRequest struct {
	FolderID int64 `json:"folder_id"`
}

func (h *FolderHandler) GetFolderContent(c echo.Context) error {
	var req getFolderContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := c.Get("user").(*store.User)
	if req.FolderID == 0 {
		root_folder_id, err := h.folderStore.GetRootFolder(user.ID)
		if err != nil {
			h.logger.Printf("Error: getting root folder id %v", err)
			return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		}
		req.FolderID = root_folder_id
	}
	folderContents, err := h.folderContentsService.GetFolderContent(user, req.FolderID)
	if err != nil {
		h.logger.Printf("Error: getting folder content %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, folderContents)
}

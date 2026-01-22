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
	logger                *log.Logger
}

func NewFolderHandler(
	folderContentsService service.FolderContentsServiceI,
	logger *log.Logger,
) *FolderHandler {
	return &FolderHandler{
		folderContentsService: folderContentsService,
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

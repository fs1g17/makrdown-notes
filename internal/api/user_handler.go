package api

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"markdown-notes/internal/store"
	"markdown-notes/internal/utils"

	"github.com/labstack/echo/v4"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore    store.UserStore
	foldersStore store.FoldersStore
	logger       *log.Logger
}

func NewUserHandler(userStore store.UserStore, foldersStore store.FoldersStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore:    userStore,
		foldersStore: foldersStore,
		logger:       logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(c echo.Context) error {
	var req registerUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	err := h.validateRegisterRequest(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("Error: hashing password %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error: creating user %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	folder_id, err := h.foldersStore.CreateFolder(user.ID, nil, "root")
	if err != nil {
		h.logger.Printf("Error: creating root folder %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, utils.Envelope{"user": user, "root": folder_id})
}

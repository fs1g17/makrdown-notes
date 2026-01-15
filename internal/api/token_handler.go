package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"markdown-notes/internal/store"
	"markdown-notes/internal/tokens"
	"markdown-notes/internal/utils"

	"github.com/labstack/echo/v4"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *createTokenRequest) validate() error {
	if len(c.Username) == 0 {
		return errors.New("username cannot be empty")
	}

	if len(c.Password) == 0 {
		return errors.New("password cannot be empty")
	}

	return nil
}

func NewTokenhandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(c echo.Context) error {
	var req createTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	if err := req.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUsername: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: PasswordHash.Match: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	if !passwordsDoMatch {
		return c.JSON(http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: Creating token: %v", err)
		return c.JSON(http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
	}

	ttl := 24 * time.Hour
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token.Plaintext // or token.String, depending on your type
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true                   // IMPORTANT: requires HTTPS
	cookie.SameSite = http.SameSiteLaxMode // or NoneMode if cross-site
	cookie.Expires = time.Now().Add(ttl)

	c.SetCookie(cookie)

	return c.JSON(http.StatusCreated, utils.Envelope{"ok": true})
}

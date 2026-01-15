package middleware

import (
	"markdown-notes/internal/store"
	"markdown-notes/internal/tokens"

	"github.com/labstack/echo/v4"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func (um *UserMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			return err
		}

		token := cookie.Value
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil || user == nil {
			return echo.NewHTTPError(401, "invalid or expired token")
		}

		c.Set("user", user)

		return next(c)
	}
}

func CurrentUser(c echo.Context) (*store.User, bool) {
	v := c.Get("user")
	if v == nil {
		return nil, false
	}
	u, ok := v.(*store.User)
	return u, ok
}

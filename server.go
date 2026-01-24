package main

import (
	"markdown-notes/internal/app"
	"markdown-notes/internal/middleware"
	"markdown-notes/internal/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	e.GET("/health", app.HealthCheck)
	e.POST("/user/register", app.UserHandler.HandleRegisterUser)
	e.POST("/tokens/auth", app.TokenHandler.HandleCreateToken)

	r := e.Group("")
	r.Use(app.UserMiddleware.AuthMiddleware)

	restricted(r, app)

	e.Logger.Fatal(e.Start(":8080"))
}

func restricted(g *echo.Group, app *app.App) {
	g.GET("/me", func(c echo.Context) error {
		u, ok := middleware.CurrentUser(c)
		if !ok {
			return echo.NewHTTPError(401, "not authenticated")
		}
		return c.JSON(200, utils.Envelope{"username": u.Username})
	})
	g.GET("/folders", app.FolderHandler.GetFolderContent)

	g.POST("/notes/new", app.NotesHandler.HandleCreateNote)
	g.POST("/folders/new", app.FolderHandler.HandleCreateFolder)
}

package app

import (
	"database/sql"
	"log"
	"markdown-notes/internal/api"
	"markdown-notes/internal/middleware"
	"markdown-notes/internal/store"
	"markdown-notes/internal/utils"
	"markdown-notes/migrations"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type App struct {
	Logger         *log.Logger
	DB             *sql.DB
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	UserMiddleware *middleware.UserMiddleware
}

func NewApp() (*App, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// our stores will go hore
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handlers will go here
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenhandler(tokenStore, userStore, logger)

	app := &App{
		Logger:       logger,
		DB:           pgDB,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		UserMiddleware: &middleware.UserMiddleware{
			UserStore: userStore,
		},
	}

	return app, nil
}

func (a *App) HealthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, utils.Envelope{"status": "ok"})
}

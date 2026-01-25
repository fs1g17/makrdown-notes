package app

import (
	"database/sql"
	"log"
	"markdown-notes/internal/api"
	"markdown-notes/internal/middleware"
	"markdown-notes/internal/service"
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
	NotesHandler   *api.NotesHandler
	FolderHandler  *api.FolderHandler
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
	notesStore := store.NewPostgresNotesStore(pgDB)
	folderStore := store.NewPostgresFoldersStore(pgDB)

	// our services will go here
	registerUserSercvice := service.NewRegisterUserService(pgDB, userStore, folderStore)
	folderContentsService := service.NewFolderContentsService(pgDB, userStore, folderStore, notesStore)

	// our handlers will go here
	userHandler := api.NewUserHandler(userStore, folderStore, registerUserSercvice, logger)
	tokenHandler := api.NewTokenhandler(tokenStore, userStore, logger)
	notesHandler := api.NewNotesHandler(notesStore, folderContentsService, logger)
	folderHandler := api.NewFolderHandler(folderContentsService, folderStore, logger)

	app := &App{
		Logger:        logger,
		DB:            pgDB,
		UserHandler:   userHandler,
		TokenHandler:  tokenHandler,
		NotesHandler:  notesHandler,
		FolderHandler: folderHandler,
		UserMiddleware: &middleware.UserMiddleware{
			UserStore: userStore,
		},
	}

	return app, nil
}

func (a *App) HealthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, utils.Envelope{"status": "ok"})
}

package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"simpletodo/internal/config"
	"simpletodo/internal/http-server/handlers"
	"simpletodo/internal/repository"
	"simpletodo/internal/service"
	"simpletodo/internal/storage/postgres"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultTimeOut      = 5 * time.Second
	defaultServeTimeOut = 15 * time.Second
)

type App struct {
	log    *zap.Logger
	router *http.Server
	db     *postgres.Storage
	cfg    config.Config
}

func NewApp(log *zap.Logger, cfg config.Config) *App {
	return &App{
		log:    log,
		cfg:    cfg,
		router: nil,
		db:     nil,
	}
}

func (a *App) Run() {
	const op = "app.Run"

	a.log.Info(op + " : starting application")

	a.db = postgres.New(a.log)
	ctx := context.Background()

	err := a.db.Init(ctx, a.cfg.Database.DatabaseURL())
	if err != nil {
		a.log.Fatal("failed to connect to the database", zap.Error(err))
	}

	tRepo := repository.NewTaskRepository(a.db, a.log)

	tService := service.New(a.log, tRepo)
	th := handlers.NewTaskHandler(tService)
	r := gin.Default()
	r = a.setupRoutes(r, th)
	a.router = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port),
		Handler:           r,
		ReadHeaderTimeout: defaultServeTimeOut,
	}

	a.log.Info(op + " : application started")

	go func() {
		a.log.Info(
			fmt.Sprintf(
				"%s : starting server on host %s and port %s",
				op,
				a.cfg.Server.Host,
				a.cfg.Server.Port,
			),
		)

		if err := a.router.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatal("failed to start server", zap.Error(err))
		}
	}()
}

func (a *App) GracefulShutdown() {
	const op = "app.GracefulShutdown"

	a.log.Info(op + " : shutting down application")

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
	defer cancel()

	if err := a.router.Shutdown(ctx); err != nil {
		a.log.Error(op+" : Server forced to shutdown", zap.Error(err))
	}

	a.log.Info(op + " : application stopped")
}

func (a *App) setupRoutes(r *gin.Engine, th *handlers.TaskHandler) *gin.Engine {
	r.LoadHTMLGlob("internal/templates/*")
	r.GET("/", th.Home)
	r.GET("/tasks", th.FetchTask)
	r.GET("/newtaskform", th.GetTaskForm)
	r.POST("/tasks", th.AddTask)
	r.GET("/gettaskupdateform/:id", th.GetTaskUpdateForm)
	r.PUT("/tasks/:id", th.UpdateTask)
	r.DELETE("/tasks/:id", th.DeleteTask)

	return r
}

package app

import (
	"context"
	"fmt"
	"net/http"
	"simpletodo/internal/config"
	"simpletodo/internal/http-server/handlers"
	"simpletodo/internal/repository"
	"simpletodo/internal/service"
	"simpletodo/internal/storage/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	log    *zap.Logger
	router *http.Server
	db     *postgres.Storage
	cfg    config.Config
}

func NewApp(log *zap.Logger, cfg config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Run() {
	const op = "app.Run"
	a.log.Info(fmt.Sprintf("%s : starting application", op))

	a.db = postgres.New(a.log)
	ctx := context.Background()
	err := a.db.Init(ctx, a.cfg.Database.DatabaseUrl())
	if err != nil {
		a.log.Fatal("failed to connect to the database", zap.Error(err))
	}

	tRepo := repository.NewTaskRepository(a.db, a.log)

	tService := service.New(a.log, tRepo)
	th := handlers.NewTaskHandler(tService)
	r := gin.Default()
	r = a.setupRoutes(r, th)
	a.router = &http.Server{Addr: fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port), Handler: r}

	a.log.Info(fmt.Sprintf("%s : application started", op))
	go func() {
		a.log.Info(fmt.Sprintf("%s : starting server on host %s and port %s", op, a.cfg.Server.Host, a.cfg.Server.Port))
		if err := a.router.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatal("failed to start server", zap.Error(err))
		}
	}()
}

func (a *App) GracefulShutdown() {
	const op = "app.GracefulShutdown"
	a.log.Info(fmt.Sprintf("%s : shutting down application", op))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.router.Shutdown(ctx); err != nil {
		a.log.Error(fmt.Sprintf("%s : Server forced to shutdown", op), zap.Error(err))
	}

	a.log.Info(fmt.Sprintf("%s : application stopped", op))
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

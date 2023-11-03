package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usersegments/config"
	v1 "usersegments/internal/controller/http/v1"
	"usersegments/internal/repository"
	"usersegments/internal/service"
	"usersegments/pkg/httpserver"
	"usersegments/pkg/logger"
	"usersegments/pkg/postgres"

	"github.com/sirupsen/logrus"
)

// @title UserSegments Service
// @version 1.0
// @description API Service for UserSegments App

// @host localhost:8080
// @BasePath /api/v1/

func Run() {
	cfg := config.NewConfig()
	logger.SetLogrus(cfg.Log.Level)

	logrus.Info("Initializing postgres...")
	db, err := postgres.New(cfg.Postgres.URL)
	if err != nil {
		logrus.Fatalf("app - Run - postgres.New. %v", err)
	}
	defer db.Close()

	logrus.Info("Initializing repositories...")
	repo := repository.NewRepository(db)

	logrus.Info("Initializing services...")
	service := service.NewService(service.Deps{
		UserRepo:      repo.UserRepo,
		SegmentRepo:   repo.SegmentRepo,
		OperationRepo: repo.OperationRepo,
	})

	logrus.Info("Initializing handlers...")
	handler := v1.NewHandler(v1.Deps{
		UserService:      service.UserService,
		OperationService: service.OperationService,
		SegmentService:   service.SegmentService,
	})

	srv := httpserver.New(cfg, handler.InitRoutes())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("error while running http server. %v\n", err.Error())
		}
	}()

	logrus.Info("Server started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		logrus.Errorf("failed to stop server. %v", err)
	}
}

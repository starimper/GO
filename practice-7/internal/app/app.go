package app

import (
	"fmt"
	"log"
	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	l := logger.New()
	l.Info("Starting practice-7 service...")

	pg, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("app.Run: postgres: %v", err)
	}
	l.Info("Connected to PostgreSQL and auto-migrated schema.")

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	router := gin.Default()
	v1.NewRouter(router, userUseCase, l)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	l.Info("Server listening on", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("app.Run: gin: %v", err)
	}
}

package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	httpDelivery "practice_3/internal/delivery/http"
	"practice_3/internal/repository"
	"practice_3/internal/repository/postgres"
	"practice_3/internal/usecase"
	"practice_3/pkg/modules"
)

func Run() {
	ctx := context.Background()

	cfg := initPostgreConfig()
	db := postgres.NewPGXDialect(ctx, cfg)

	repos := repository.NewRepositories(db)
	userUsecase := usecase.NewUserUsecase(repos.UserRepository)
	handler := httpDelivery.NewHandler(userUsecase)

	router := mux.NewRouter()
	handler.Register(router)

	router.Use(httpDelivery.LoggingMiddleware)
	router.Use(httpDelivery.AuthMiddleware("secret123"))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "1234",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	httpDelivery "practice_4/internal/delivery/http"
	"practice_4/internal/repository"
	"practice_4/internal/repository/postgres"
	"practice_4/internal/usecase"
	"practice_4/pkg/modules"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := initPostgreConfigFromEnv()

	log.Println("Waiting for database to be ready...")
	db := postgres.NewPGXDialect(ctx, cfg) // внутри есть Ping

	repos := repository.NewRepositories(db)
	userUsecase := usecase.NewUserUsecase(repos.UserRepository)
	handler := httpDelivery.NewHandler(userUsecase)

	router := mux.NewRouter()
	handler.Register(router)

	router.Use(httpDelivery.LoggingMiddleware)
	router.Use(httpDelivery.AuthMiddleware(os.Getenv("API_KEY")))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Starting the Server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	_ = db.DB.Close()
}

func initPostgreConfigFromEnv() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Username:    os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
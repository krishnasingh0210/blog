package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/gorilla/mux"

	"github.com/krishnasingh0210/blog/internal/api/handlers"
	"github.com/krishnasingh0210/blog/internal/api/routes"
	"github.com/krishnasingh0210/blog/internal/infrastructure"
	"github.com/krishnasingh0210/blog/internal/repository"
	"github.com/krishnasingh0210/blog/internal/service"
	"github.com/krishnasingh0210/blog/pkg/jwt"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtManager)
	authHandler := handlers.NewAuthHandler(userService)

	r := mux.NewRouter()
	routes.Register(r, routes.Handlers{Auth: authHandler}, jwtManager)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	log.Printf("server starting on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

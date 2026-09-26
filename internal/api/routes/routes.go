package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/krishnasingh0210/blog/internal/api/handlers"
	"github.com/krishnasingh0210/blog/internal/middleware"
	"github.com/krishnasingh0210/blog/pkg/jwt"
)

type Handlers struct {
	Auth *handlers.AuthHandler
}

func Register(r *mux.Router, h Handlers, jwtManager *jwt.Manager) {
	// public routes — no auth required
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/auth/register", h.Auth.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/login", h.Auth.Login).Methods(http.MethodPost)

	// example of a protected subrouter for later phases (blogs, comments):
	// authed := r.PathPrefix("/api/v1").Subrouter()
	// authed.Use(middleware.Auth(jwtManager))
	// adminOnly := authed.PathPrefix("/admin").Subrouter()
	// adminOnly.Use(middleware.RequireRole("admin"))
	_ = middleware.Auth // silence unused import until Phase 5 adds blog routes
}
package main

import (
	"database/sql"
	"log"
	"net/http"
	"profiles"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

const AuthHashCost = 8

func main() {
	// profiles
	sql, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalf("error while opening sql db: %v", err)
	}
	onNewRegister := profiles.NewRegisterCallback(sql)
	profilesRouter := profiles.NewProfilesRouterImpl(sql)

	// auth
	authStore, err := auth.NewStoreImpl("auth.db.csv")
	if err != nil {
		log.Fatalf("error while opening auth store: %v", err)
	}
	loginHandler, registerHandler := auth.NewHandlersImpl(authStore, AuthHashCost, onNewRegister)
	authMiddleware := auth.NewTokenAuthMiddleware(authStore).Middleware

	// routing
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", loginHandler.ServeHTTP)
		r.Post("/register", registerHandler.ServeHTTP)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Route("/profiles", profilesRouter)
	})

	http.ListenAndServe(":4242", r)
}

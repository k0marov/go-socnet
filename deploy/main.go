package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/features/comments"
	"github.com/k0marov/socnet/features/posts"
	"github.com/k0marov/socnet/features/profiles"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	auth "github.com/k0marov/golang-auth"
)

const AuthHashCost = 8

func main() {
	// db
	sql, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalf("error while opening sql db: %v", err)
	}

	// profiles
	onNewRegister := profiles.NewRegisterCallback(sql)
	profileGetter := profiles.NewProfileGetterImpl(sql)
	profilesRouter := profiles.NewProfilesRouterImpl(sql)

	// posts
	postsRouter := posts.NewPostsRouterImpl(sql, profileGetter)

	// comments
	commentsRouter := comments.NewCommentsRouterImpl(sql, profileGetter)

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
		r.Route("/posts", postsRouter)
		r.Route("/comments", commentsRouter)
	})

	http.ListenAndServe(":4242", r)
}

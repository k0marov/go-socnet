package core

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/general/periodic"
	"github.com/k0marov/go-socnet/features/comments"
	"github.com/k0marov/go-socnet/features/feed"
	"github.com/k0marov/go-socnet/features/posts"
	"github.com/k0marov/go-socnet/features/profiles"
	auth "github.com/k0marov/golang-auth"
	"log"
	"net/http"
	"time"
)

const AuthHashCost = 8

func Setup() http.Handler {
	sql, err := sqlx.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatalf("error while opening sql db: %v", err)
	}
	sql.Exec("PRAGMA foreign_keys = ON;")

	// profiles
	onNewRegister := profiles.NewRegisterCallback(sql)
	profileGetter := profiles.NewProfileGetterImpl(sql)
	profilesRouter := profiles.NewProfilesRouterImpl(sql)

	// posts
	postsRouter := posts.NewPostsRouterImpl(sql, profileGetter)
	postRecommendable := posts.NewPostRecommendable(sql)
	periodic.RunPeriodically(func() {
		err := postRecommendable.UpdateRecs()
		if err != nil {
			log.Printf("while updating recommendations for posts: %v", err)
		}
	}, 1*time.Minute)

	// feed
	feedRouter := feed.NewFeedRouterImpl(sql, postRecommendable)

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
		r.Route("/feed", feedRouter)
	})

	return r
}

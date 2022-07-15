package feed

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/core/abstract/recommendable"
	"github.com/k0marov/go-socnet/features/feed/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/feed/delivery/http/router"
	"github.com/k0marov/go-socnet/features/feed/domain/service"
)

func NewFeedRouterImpl(db *sql.DB, postRecommendable recommendable.Recommendable) func(chi.Router) {
	// service
	getFeed := service.NewFeedGetter(postRecommendable.GetRecs)
	// handlers
	feedHandler := handlers.NewFeedHandler(getFeed)

	return router.NewFeedRouter(feedHandler)
}

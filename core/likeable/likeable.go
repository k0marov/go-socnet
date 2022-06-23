package likeable

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/likeable/service"
	"github.com/k0marov/socnet/core/likeable/store/sql_db"
	"github.com/k0marov/socnet/core/likeable/table_name"
)

type likeable struct {
	ToggleLike    service.LikeToggler
	GetLikesCount service.LikesCountGetter
	IsLiked       service.LikeChecker
}

func NewLikeable(db *sql.DB, targetTableName string) (likeable, error) {
	// store
	store, err := sql_db.NewSqlDB(db, table_name.NewTableName(targetTableName))
	if err != nil {
		return likeable{}, fmt.Errorf("while opening the likeable sql db: %w", err)
	}
	// service
	toggleLike := service.NewLikeToggler(store.IsLiked, store.Like, store.Unlike)
	getLikesCount := service.NewLikesCountGetter(store.GetLikesCount)
	isLiked := service.NewLikeChecker(store.IsLiked)
	return likeable{
		ToggleLike:    toggleLike,
		GetLikesCount: getLikesCount,
		IsLiked:       isLiked,
	}, nil
}

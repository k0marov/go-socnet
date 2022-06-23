package likeable

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable/service"
	"github.com/k0marov/socnet/core/likeable/store/sql_db"
	"github.com/k0marov/socnet/core/likeable/table_name"
)

type LikeToggler func(id string, owner, liker core_values.UserId) error
type LikesCountGetter func(id string) (int, error)
type LikeChecker func(id string, fromUser core_values.UserId) (bool, error)

type likeable struct {
	ToggleLike    LikeToggler
	GetLikesCount LikesCountGetter
	IsLiked       LikeChecker
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

package likeable

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/likeable/service"
	"github.com/k0marov/socnet/core/likeable/store/sql_db"
	"github.com/k0marov/socnet/core/likeable/table_name"
)

type (
	LikeToggler          = service.LikeToggler
	LikeChecker          = service.LikeChecker
	LikesCountGetter     = service.LikesCountGetter
	UserLikesCountGetter = service.UserLikesCountGetter
	UserLikesGetter      = service.UserLikesGetter
)

type likeable struct {
	ToggleLike        LikeToggler
	IsLiked           LikeChecker
	GetLikesCount     LikesCountGetter
	GetUserLikesCount UserLikesCountGetter
	GetUserLikes      UserLikesGetter
}

func NewLikeable(db *sql.DB, targetTableName table_name.TableName) (likeable, error) {
	// store
	store, err := sql_db.NewSqlDB(db, targetTableName)
	if err != nil {
		return likeable{}, fmt.Errorf("while opening the likeable sql db: %w", err)
	}
	// service
	toggleLike := service.NewLikeToggler(store.IsLiked, store.Like, store.Unlike)
	isLiked := service.NewLikeChecker(store.IsLiked)
	getLikesCount := service.NewLikesCountGetter(store.GetLikesCount)
	getUserLikesCount := service.NewUserLikesCountGetter(store.GetUserLikesCount)
	getUserLikes := service.NewUserLikesGetter(store.GetUserLikes)
	return likeable{
		ToggleLike:        toggleLike,
		IsLiked:           isLiked,
		GetLikesCount:     getLikesCount,
		GetUserLikesCount: getUserLikesCount,
		GetUserLikes:      getUserLikes,
	}, nil
}

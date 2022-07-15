package likeable

import (
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/likeable/service"
	"github.com/k0marov/go-socnet/core/abstract/likeable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
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

func NewLikeable(db *sqlx.DB, targetTableName table_name.TableName) (likeable, error) {
	// store
	store, err := sql_db.NewSqlDB(db, targetTableName)
	if err != nil {
		return likeable{}, core_err.Rethrow("opening the likeable sql db", err)
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

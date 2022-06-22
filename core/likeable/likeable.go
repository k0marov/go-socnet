package likeable

import (
	"database/sql"
	"github.com/k0marov/socnet/core/core_values"
)

type LikeToggler func(id string, owner, liker core_values.UserId) error
type LikesCountGetter func(id string) (int, error)
type LikeChecker func(id string, fromUser core_values.UserId) (bool, error)

type likeable struct {
	ToggleLike    LikeToggler
	GetLikesCount LikesCountGetter
	IsLiked       LikeChecker
}

func NewLikeable(db *sql.DB, targetTableName string) likeable {
	return likeable{}
}

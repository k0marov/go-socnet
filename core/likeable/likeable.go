package likeable

import (
	"database/sql"
	"github.com/k0marov/socnet/core/core_values"
)

type Likeable interface {
	LikeToggler
	LikesCountGetter
	LikeChecker
}

type LikeToggler func(id string, fromUser core_values.UserId) error
type LikesCountGetter func(id string) (int, error)
type LikeChecker func(id string, fromUser core_values.UserId) (bool, error)

type likeable struct {
	ToggleLike    LikeToggler
	GetLikesCount LikesCountGetter
	IsLiked       LikeChecker
}

func NewLikeable(db *sql.DB, entityName string) likeable {
	return likeable{}
}

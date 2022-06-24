package contexters

import (
	"fmt"

	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/likeable"
)

type LikeableContext struct {
	IsLiked bool
	IsMine  bool
}

type LikeableContextGetter func(target string, owner, caller core_values.UserId) (LikeableContext, error)

func NewLikeableContextGetter(checkLiked likeable.LikeChecker) LikeableContextGetter {
	return func(target string, owner, caller core_values.UserId) (LikeableContext, error) {
		isLiked, err := checkLiked(target, caller)
		if err != nil {
			return LikeableContext{}, fmt.Errorf("while checking if target is liked: %w", err)
		}
		return LikeableContext{
			IsLiked: isLiked,
			IsMine:  owner == caller,
		}, nil
	}
}

package contexters

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type OwnLikeContext struct {
	IsLiked bool
	IsMine  bool
}

type OwnLikeContextGetter func(target string, owner, caller core_values.UserId) (OwnLikeContext, error)

func NewOwnLikeContextGetter(checkLiked likeable.LikeChecker) OwnLikeContextGetter {
	return func(target string, owner, caller core_values.UserId) (OwnLikeContext, error) {
		isLiked, err := checkLiked(target, caller)
		if err != nil {
			return OwnLikeContext{}, fmt.Errorf("while checking if target is liked: %w", err)
		}
		return OwnLikeContext{
			IsLiked: isLiked,
			IsMine:  caller == owner,
		}, nil
	}
}

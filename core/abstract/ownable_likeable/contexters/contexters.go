package contexters

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type OwnLikeContext struct {
	IsLiked bool
	IsMine  bool
}

type OwnLikeContextGetter func(target string, caller core_values.UserId) (OwnLikeContext, error)

func NewOwnLikeContextGetter(checkLiked likeable.LikeChecker, getOwner ownable.OwnerGetter) OwnLikeContextGetter {
	return func(target string, caller core_values.UserId) (OwnLikeContext, error) {
		isLiked, err := checkLiked(target, caller)
		if err != nil {
			return OwnLikeContext{}, fmt.Errorf("while checking if target is liked: %w", err)
		}
		owner, err := getOwner(target)
		if err != nil {
			return OwnLikeContext{}, fmt.Errorf("while getting target's owner: %w", err)
		}
		return OwnLikeContext{
			IsLiked: isLiked,
			IsMine:  caller == owner,
		}, nil
	}
}

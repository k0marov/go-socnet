package contexters

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/helpers"
	"github.com/k0marov/socnet/core/likeable"
)

type LikeableContext struct {
	IsLiked bool
	IsMine  bool
}

type LikeableTarget interface {
	GetId() string
	GetOwner() core_values.UserId
}

type LikeableContextGetter func(target LikeableTarget, caller core_values.UserId) (LikeableContext, error)
type LikeableListContextsGetter func(targets []LikeableTarget, caller core_values.UserId) ([]LikeableContext, error)

func NewLikeableContextGetter(checkLiked likeable.LikeChecker) LikeableContextGetter {
	return func(target LikeableTarget, caller core_values.UserId) (LikeableContext, error) {
		isLiked, err := checkLiked(target.GetId(), caller)
		if err != nil {
			return LikeableContext{}, fmt.Errorf("while checking if target is liked: %w", err)
		}
		return LikeableContext{
			IsLiked: isLiked,
			IsMine:  target.GetOwner() == caller,
		}, nil
	}
}

func NewLikeableListContextsGetter(getContext LikeableContextGetter) LikeableListContextsGetter {
	return func(targets []LikeableTarget, caller core_values.UserId) ([]LikeableContext, error) {
		return helpers.MapForEach(targets, func(target LikeableTarget) (LikeableContext, error) {
			return getContext(target, caller)
		})
	}
}

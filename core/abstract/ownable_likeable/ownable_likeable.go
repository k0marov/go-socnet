package ownable_likeable

import (
	"github.com/k0marov/go-socnet/core/abstract/likeable"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/ownable_likeable/service"
)

type SafeLikeToggler = service.SafeLikeToggler

type ownableLikeable struct {
	SafeToggleLike SafeLikeToggler
}

func NewOwnableLikeable(getOwner ownable.OwnerGetter, toggleLike likeable.LikeToggler) ownableLikeable {
	safeToggleLike := service.NewSafeLikeToggler(getOwner, toggleLike)
	return ownableLikeable{SafeToggleLike: safeToggleLike}
}

package store_contracts

import (
	"core/core_values"
	"posts/domain/entities"
)

type StorePostsGetter func(authorId core_values.UserId) ([]entities.Post, error)

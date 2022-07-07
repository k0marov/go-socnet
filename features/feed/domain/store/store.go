package store

import "github.com/k0marov/go-socnet/features/posts/domain/entities"

type RandomPostsGetter = func(count int) ([]entities.Post, error)

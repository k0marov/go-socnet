package contexters

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/entities"
)

type PostContextAdder func(post entities.Post, caller core_values.UserId) (entities.ContextedPost, error)
type PostListContextAdder func(posts []entities.Post, caller core_values.UserId) ([]entities.ContextedPost, error)

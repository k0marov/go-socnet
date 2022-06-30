package contexters

import (
	"fmt"

	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/helpers"
	likeable_contexters "github.com/k0marov/go-socnet/core/likeable/contexters"
	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
)

type PostContextAdder func(post entities.Post, caller core_values.UserId) (entities.ContextedPost, error)
type PostListContextAdder func(posts []entities.Post, caller core_values.UserId) ([]entities.ContextedPost, error)

func NewPostContextAdder(getProfile profile_service.ProfileGetter, getContext likeable_contexters.LikeableContextGetter) PostContextAdder {
	return func(post entities.Post, caller core_values.UserId) (entities.ContextedPost, error) {
		author, err := getProfile(post.PostModel.AuthorId, caller)
		if err != nil {
			return entities.ContextedPost{}, fmt.Errorf("while getting author of post: %w", err)
		}
		context, err := getContext(post.Id, author.Id, caller)
		if err != nil {
			return entities.ContextedPost{}, fmt.Errorf("while getting context of post: %w", err)
		}
		ctxPost := entities.ContextedPost{
			Post:            post,
			LikeableContext: context,
			Author:          author,
		}
		return ctxPost, nil
	}
}

func NewPostListContextAdder(addContext PostContextAdder) PostListContextAdder {
	return func(posts []entities.Post, caller core_values.UserId) ([]entities.ContextedPost, error) {
		return helpers.MapForEachWithErr(posts, func(post entities.Post) (entities.ContextedPost, error) {
			return addContext(post, caller)
		})
	}
}

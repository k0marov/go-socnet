package contexters

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"

	"github.com/k0marov/go-socnet/core/helpers"
	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
)

type PostContextAdder func(post entities.Post, caller core_values.UserId) (entities.ContextedPost, error)
type PostListContextAdder func(posts []entities.Post, caller core_values.UserId) ([]entities.ContextedPost, error)

func NewPostContextAdder(getProfile profile_service.ProfileGetter, getContext likeable_contexters.OwnLikeContextGetter) PostContextAdder {
	return func(post entities.Post, caller core_values.UserId) (entities.ContextedPost, error) {
		author, err := getProfile(post.PostModel.AuthorId, caller)
		if err != nil {
			return entities.ContextedPost{}, core_err.Rethrow("getting author of post", err)
		}
		context, err := getContext(post.Id, author.Id, caller)
		if err != nil {
			return entities.ContextedPost{}, core_err.Rethrow("getting context of post", err)
		}
		ctxPost := entities.ContextedPost{
			Post:           post,
			OwnLikeContext: context,
			Author:         author,
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

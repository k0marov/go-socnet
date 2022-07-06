package contexters

import (
	likeable_contexters "github.com/k0marov/go-socnet/core/abstract/ownable_likeable/contexters"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"

	"github.com/k0marov/go-socnet/core/helpers"
	"github.com/k0marov/go-socnet/features/comments/domain/entities"
	profile_service "github.com/k0marov/go-socnet/features/profiles/domain/service"
)

type CommentContextAdder func(comment entities.Comment, caller core_values.UserId) (entities.ContextedComment, error)
type CommentListContextAdder func(comments []entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error)

func NewCommentContextAdder(getProfile profile_service.ProfileGetter, getContext likeable_contexters.OwnLikeContextGetter) CommentContextAdder {
	return func(comment entities.Comment, caller core_values.UserId) (entities.ContextedComment, error) {
		author, err := getProfile(comment.AuthorId, caller)
		if err != nil {
			return entities.ContextedComment{}, core_err.Rethrow("getting author of comment", err)
		}
		context, err := getContext(comment.Id, author.Id, caller)
		if err != nil {
			return entities.ContextedComment{}, core_err.Rethrow("getting context data of a comment", err)
		}
		return entities.ContextedComment{
			Comment:        comment,
			Author:         author,
			OwnLikeContext: context,
		}, nil
	}
}

func NewCommentListContextAdder(addContext CommentContextAdder) CommentListContextAdder {
	return func(comments []entities.Comment, caller core_values.UserId) (ctxComments []entities.ContextedComment, err error) {
		return helpers.MapForEachWithErr(comments, func(comm entities.Comment) (entities.ContextedComment, error) {
			return addContext(comm, caller)
		})
	}
}

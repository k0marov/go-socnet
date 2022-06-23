package contexters

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/helpers"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	profile_service "github.com/k0marov/socnet/features/profiles/domain/service"
)

type CommentContextAdder func(comment entities.Comment, caller core_values.UserId) (entities.ContextedComment, error)
type CommentListContextAdder func(comments []entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error)

func NewCommentContextAdder(getProfile profile_service.ProfileGetter, getContext likeable_contexters.LikeableContextGetter) CommentContextAdder {
	return func(comment entities.Comment, caller core_values.UserId) (entities.ContextedComment, error) {
		author, err := getProfile(comment.AuthorId, caller)
		if err != nil {
			return entities.ContextedComment{}, fmt.Errorf("while getting author of comment: %w", err)
		}
		context, err := getContext(comment.Id, author.Id, caller)
		if err != nil {
			return entities.ContextedComment{}, fmt.Errorf("while getting context data of a comment: %w", err)
		}
		return entities.ContextedComment{
			Comment:         comment,
			Author:          author,
			LikeableContext: context,
		}, nil
	}
}

func NewCommentListContextAdder(addContext CommentContextAdder) CommentListContextAdder {
	return func(comments []entities.Comment, caller core_values.UserId) (ctxComments []entities.ContextedComment, err error) {
		return helpers.MapForEach(comments, func(comm entities.Comment) (entities.ContextedComment, error) {
			return addContext(comm, caller)
		})
	}
}

package contexters

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/entities"
)

type CommentContextAdder func(comment entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error)
type CommentListContextAdder func(comments []entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error)

func NewCommentContextAdder() CommentContextAdder {
	return func(comment entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error) {
		panic("unimplemented")
	}
}

func NewCommentListContextAdder() CommentListContextAdder {
	return func(comments []entities.Comment, caller core_values.UserId) ([]entities.ContextedComment, error) {
		panic("unimplemented")
	}
}

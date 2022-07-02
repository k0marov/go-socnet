package validators

import (
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
)

type CommentValidator func(values.NewCommentValue) (client_errors.ClientError, bool)

const MaxTextLength = 255

func NewCommentValidator() CommentValidator {
	return func(newComment values.NewCommentValue) (client_errors.ClientError, bool) {
		if newComment.Text == "" {
			return client_errors.EmptyText, false
		}
		if len(newComment.Text) > MaxTextLength {
			return client_errors.TextTooLong, false
		}
		return client_errors.ClientError{}, true
	}
}

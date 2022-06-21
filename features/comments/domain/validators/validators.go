package validators

import (
	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/features/comments/domain/values"
)

type CommentValidator func(values.NewCommentValue) (client_errors.ClientError, bool)

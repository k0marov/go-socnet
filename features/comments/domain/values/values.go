package values

import (
	"github.com/k0marov/go-socnet/core/core_values"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
)

type CommentId = string

type NewCommentValue struct {
	Author core_values.UserId
	Post   post_values.PostId
	Text   string
}

package store

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
)

type (
	CommentsGetter func(post post_values.PostId) ([]models.CommentModel, error)
	LikeChecker    func(comment values.CommentId, caller core_values.UserId) (bool, error)
	Liker          func(comment values.CommentId, liker core_values.UserId) error
	Unliker        func(comment values.CommentId, unliker core_values.UserId) error
	Creator        func(newComment values.NewCommentValue) (models.CommentModel, error)
)

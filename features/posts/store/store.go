package store

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/store"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"time"
)

type (
	DBPostsGetter  func(core_values.UserId) ([]entities.Post, error)
	DBLiker        func(values.PostId, core_values.UserId) error
	DBUnliker      func(values.PostId, core_values.UserId) error
	DBLikeChecker  func(values.PostId, core_values.UserId) (bool, error)
	DBAuthorGetter func(values.PostId) (core_values.UserId, error)
)

func NewStorePostCreator() store.StorePostCreator {
	return func(post values.NewPostData, createdAt time.Time) error {
		panic("unimplemented")
	}
}

func NewStorePostDeleter() store.StorePostDeleter {
	return func(postId values.PostId) error {
		panic("unimplemented")
	}
}

func NewStorePostsGetter(getter DBPostsGetter) store.StorePostsGetter {
	return store.StorePostsGetter(getter)
}

func NewStoreLiker(liker DBLiker) store.StoreLiker {
	return store.StoreLiker(liker)
}

func NewStoreLikeChecker(likeChecker DBLikeChecker) store.StoreLikeChecker {
	return store.StoreLikeChecker(likeChecker)
}
func NewStoreUnliker(unliker DBUnliker) store.StoreUnliker {
	return store.StoreUnliker(unliker)
}

func NewStoreAuthorGetter(authorGetter DBAuthorGetter) store.StoreAuthorGetter {
	return store.StoreAuthorGetter(authorGetter)
}

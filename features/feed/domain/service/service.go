package service

import (
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/features/feed/domain/store"
	"github.com/k0marov/go-socnet/features/posts/domain/contexters"
	"github.com/k0marov/go-socnet/features/posts/domain/entities"
	"strconv"
)

type FeedGetter = func(count string, caller core_values.UserId) ([]entities.ContextedPost, error)

const DefaultCount = 5
const MaxCount = 50

func NewFakeFeedGetter(getPosts store.RandomPostsGetter, addContext contexters.PostListContextAdder) FeedGetter {
	return func(countStr string, caller core_values.UserId) ([]entities.ContextedPost, error) {
		var count int
		if countStr == "" {
			count = DefaultCount
		} else {
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil {
				return []entities.ContextedPost{}, client_errors.NonIntegerCount
			}
		}
		if count > MaxCount {
			return []entities.ContextedPost{}, client_errors.TooBigCount
		}
		posts, err := getPosts(count)
		if err != nil {
			return []entities.ContextedPost{}, core_err.Rethrow("getting posts", err)
		}
		postsWithCtx, err := addContext(posts, caller)
		return postsWithCtx, nil
	}
}

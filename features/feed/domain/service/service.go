package service

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/features/posts/domain/contexters"
)

type FeedGetter = func(count string, caller core_values.UserId) ([]string, error)

const DefaultCount = 5
const MaxCount = 50

func NewFakeFeedGetter(addContext contexters.PostListContextAdder) FeedGetter {
	return func(countStr string, caller core_values.UserId) ([]string, error) {
		panic("unimplemented")
		//var count int
		//if countStr == "" {
		//	count = DefaultCount
		//} else {
		//	var err error
		//	count, err = strconv.Atoi(countStr)
		//	if err != nil {
		//		return []entities.ContextedPost{}, client_errors.NonIntegerCount
		//	}
		//}
		//if count > MaxCount {
		//	return []entities.ContextedPost{}, client_errors.TooBigCount
		//}
		//posts, err := getPosts(count)
		//if err != nil {
		//	return []entities.ContextedPost{}, core_err.Rethrow("getting posts", err)
		//}
		//postsWithCtx, err := addContext(posts, caller)
		//return postsWithCtx, nil
	}
}

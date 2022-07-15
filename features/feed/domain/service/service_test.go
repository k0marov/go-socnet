package service_test

import (
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"github.com/k0marov/go-socnet/features/feed/domain/service"
	"testing"
)

func TestFeedGetter(t *testing.T) {
	caller := RandomId()
	countStr := "8"
	posts := []string{RandomId(), RandomId(), RandomId()}

	t.Run("error case - count is not int", func(t *testing.T) {
		_, err := service.NewFeedGetter(nil)("asdf", caller)
		AssertError(t, err, client_errors.NonIntegerCount)
	})
	t.Run("error case - count is too big", func(t *testing.T) {
		_, err := service.NewFeedGetter(nil)("9999", caller)
		AssertError(t, err, client_errors.TooBigCount)
	})

	feedGetter := func(callerId core_values.UserId, count int) ([]string, error) {
		if count == 8 && callerId == caller {
			return posts, nil
		}
		panic("unexpected")
	}
	t.Run("error case - getting feed throws", func(t *testing.T) {
		feedGetter := func(core_values.UserId, int) ([]string, error) {
			return nil, RandomError()
		}
		_, err := service.NewFeedGetter(feedGetter)(countStr, caller)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		gotPosts, err := service.NewFeedGetter(feedGetter)(countStr, caller)
		AssertNoError(t, err)
		Assert(t, gotPosts, posts, "returned posts")
	})
	t.Run("happy case - count is empty", func(t *testing.T) {
		feedGetter := func(callerId core_values.UserId, count int) ([]string, error) {
			if count == service.DefaultCount && callerId == caller {
				return posts, nil
			}
			panic("unexpected")
		}
		gotPosts, err := service.NewFeedGetter(feedGetter)("", caller)
		AssertNoError(t, err)
		Assert(t, gotPosts, posts, "returned posts")
	})

}

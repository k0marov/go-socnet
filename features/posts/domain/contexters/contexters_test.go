package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/domain/contexters"
	"github.com/k0marov/socnet/features/posts/domain/entities"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"testing"
)

func TestPostContextAdder(t *testing.T) {
	post := RandomPost()
	caller := RandomId()
	author := RandomContextedProfile()
	ctx := RandomLikeableContext()

	getProfile := func(id, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
		if id == post.PostModel.AuthorId && callerId == caller {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting author throws", func(t *testing.T) {
		getProfile := func(id, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
			return profile_entities.ContextedProfile{}, RandomError()
		}
		_, err := contexters.NewPostContextAdder(getProfile, nil)(post, caller)
		AssertSomeError(t, err)
	})
	getContext := func(target string, owner, callerId core_values.UserId) (likeable_contexters.LikeableContext, error) {
		if target == post.Id && owner == author.Id && callerId == caller {
			return ctx, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting context throws", func(t *testing.T) {
		getContext := func(string, core_values.UserId, core_values.UserId) (likeable_contexters.LikeableContext, error) {
			return likeable_contexters.LikeableContext{}, RandomError()
		}
		_, err := contexters.NewPostContextAdder(getProfile, getContext)(post, caller)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		wantPost := entities.ContextedPost{
			Post:            post,
			Author:          author,
			LikeableContext: ctx,
		}
		gotPost, err := contexters.NewPostContextAdder(getProfile, getContext)(post, caller)
		AssertNoError(t, err)
		Assert(t, gotPost, wantPost, "returned post")
	})
}

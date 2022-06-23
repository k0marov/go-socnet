package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	likeable_contexters "github.com/k0marov/socnet/core/likeable/contexters"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/contexters"
	"github.com/k0marov/socnet/features/comments/domain/entities"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"testing"
)

func TestCommentContextAdder(t *testing.T) {
	comment := RandomComment()
	author := RandomContextedProfile()
	caller := RandomId()
	context := likeable_contexters.LikeableContext{
		IsLiked: RandomBool(),
		IsMine:  RandomBool(),
	}

	authorGetter := func(id, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
		if id == comment.AuthorId && callerId == caller {
			return author, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting author throws", func(t *testing.T) {
		authorGetter := func(id, callerId core_values.UserId) (profile_entities.ContextedProfile, error) {
			return profile_entities.ContextedProfile{}, RandomError()
		}
		_, err := contexters.NewCommentContextAdder(authorGetter, nil)(comment, caller)
		AssertSomeError(t, err)
	})
	contextGetter := func(target string, ownerId, callerId core_values.UserId) (likeable_contexters.LikeableContext, error) {
		if target == comment.Id && ownerId == author.Id && callerId == caller {
			return context, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting likeable context throws", func(t *testing.T) {
		contextGetter := func(target string, ownerId, callerId core_values.UserId) (likeable_contexters.LikeableContext, error) {
			return likeable_contexters.LikeableContext{}, RandomError()
		}
		_, err := contexters.NewCommentContextAdder(authorGetter, contextGetter)(comment, caller)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		contextedComment, err := contexters.NewCommentContextAdder(authorGetter, contextGetter)(comment, caller)
		AssertNoError(t, err)
		wantComment := entities.ContextedComment{
			Comment:         comment,
			LikeableContext: context,
			Author:          author,
		}
		Assert(t, contextedComment, wantComment, "returned contexted comment")
	})
}

package contexters_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/contexters"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"testing"
)

func TestCommentContextAdder(t *testing.T) {
	comment := RandomComment()
	author := RandomContextedProfile()
	isLiked := RandomBool()
	isMine := RandomBool()
	var caller core_values.UserId
	if isMine {
		caller = author.Id
	} else {
		caller = RandomId()
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
	checkLiked := func(targetId string, callerId core_values.UserId) (bool, error) {
		if targetId == comment.Id && callerId == caller {
			return isLiked, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - checking if comment is liked throws", func(t *testing.T) {
		checkLiked := func(string, core_values.UserId) (bool, error) {
			return false, RandomError()
		}
		_, err := contexters.NewCommentContextAdder(authorGetter, checkLiked)(comment, caller)
		AssertSomeError(t, err)
	})
	contextedComment, err := contexters.NewCommentContextAdder(authorGetter, checkLiked)(comment, caller)
	AssertNoError(t, err)
	Assert(t, contextedComment.Comment, comment, "the contexted comment's comment value")
	Assert(t, contextedComment.IsLiked, isLiked, "isLiked")
	Assert(t, contextedComment.IsMine, isMine, "isMine")
}

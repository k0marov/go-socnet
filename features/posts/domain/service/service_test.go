package service_test

import (
	"core/core_values"
	. "core/test_helpers"
	"posts/domain/entities"
	"posts/domain/service"
	"testing"
)

func TestPostsGetter(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		randomPosts := []entities.Post{RandomPost(), RandomPost(), RandomPost()}
		randomAuthor := RandomString()
		storePostsGetter := func(authorId core_values.UserId) ([]entities.Post, error) {
			if authorId == randomAuthor {
				return randomPosts, nil
			}
			panic("unexpected args")
		}
		sut := service.NewPostsGetter(storePostsGetter)
		gotPosts, err := sut(randomAuthor)
		AssertNoError(t, err)
		Assert(t, gotPosts, randomPosts, "the returned posts")
	})
	t.Run("error case - store throws an error", func(t *testing.T) {
		storeGetter := func(core_values.UserId) ([]entities.Post, error) {
			return []entities.Post{}, RandomError()
		}
		_, err := service.NewPostsGetter(storeGetter)("42")
		AssertSomeError(t, err)
	})
}

func TestPostDeleter(t *testing.T) {
	t.Run("happy case", func(t *testing.T) {
		// authorGetter := func(postId values.PostId) (core_values.UserId, error) {

		// }
	})
	t.Run("error case - the calling user is not the post author", func(t *testing.T) {

	})
	t.Run("error case - store returns 404", func(t *testing.T) {

	})
	t.Run("error case - store returns some other error", func(t *testing.T) {

	})
}

func TestPostCreator(t *testing.T) {

}

func TestPostLikeToggler(t *testing.T) {

}

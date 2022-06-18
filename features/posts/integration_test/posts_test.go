package posts_test

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts"
	"github.com/k0marov/socnet/features/posts/domain/entities"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/profiles"
	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"os"
	"path/filepath"
	"testing"
)

func TestPosts(t *testing.T) {
	// working directory
	os.Mkdir("tmp_test", 0777)
	os.Chdir("tmp_test")
	defer func() {
		os.Chdir("..")
		os.RemoveAll("tmp_test")
	}()

	// db
	sql, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error while opening in-memory database: %v", err)
	}

	r := chi.NewRouter()
	// profiles
	r.Route("/profiles", profiles.NewProfilesRouterImpl(sql))
	fakeRegisterProfile := profiles.NewRegisterCallback(sql)
	// posts
	r.Route("/posts", posts.NewPostsRouterImpl(sql, profiles.NewProfileGetterImpl(sql)))

	// helpers
	//addAuthToReq := func(req *http.Request, user core_entities.User) *http.Request {
	//	ctx := req.Context()
	//	ctx = context.WithValue(ctx, auth.UserContextKey, auth.User{Id: user.Id, Username: user.Username})
	//	return req.WithContext(ctx)
	//}
	createPost := func(t testing.TB, author auth.User, images [][]byte, text string) {
		t.Helper()
		panic("unimplemented")
	}
	getPosts := func(t testing.TB, author core_values.UserId, caller auth.User) []entities.ContextedPost {
		t.Helper()
		panic("unimplemented")
	}
	assertImageCreated := func(t testing.TB, postImage values.PostImage, wantImage []byte) {
		t.Helper()
		panic("unimplemented")
	}
	deletePost := func(t testing.TB, postId values.PostId, author auth.User) {
		t.Helper()
		panic("unimplemented")
	}
	assertPostFilesDeleted := func(t testing.TB, postId values.PostId, author core_values.UserId) {
		t.Helper()
		panic("unimplemented")
	}
	toggleLike := func(t testing.TB, postId values.PostId, caller auth.User) {
		t.Helper()
		panic("unimplemented")
	}

	registerProfile := func(user auth.User) profile_entities.Profile {
		fakeRegisterProfile(user)
		return profile_entities.Profile{
			Id:       user.Id,
			Username: user.Username,
		}
	}

	// create 2 profiles
	user1 := RandomAuthUser()
	user2 := RandomAuthUser()
	registerProfile(user1)
	registerProfile(user2)

	t.Run("creating, reading and deleting posts", func(t *testing.T) {
		// create 1 post (with images) belonging to 2-nd profile
		text2 := "Hello, World with Images!"
		image1 := readFixture(t, "test_image.jpg")
		image2 := readFixture(t, "test_image.jpg")
		createPost(t, user2, [][]byte{image1, image2}, text2)

		// create 1 post (without images) belonging to 2-nd profile
		text1 := "Hello, World!"
		createPost(t, user2, [][]byte{}, text1)

		// assert they were created
		var posts []entities.ContextedPost
		posts = getPosts(t, user2.Id, user2)

		Assert(t, len(posts), 2, "number of created posts")

		Assert(t, posts[0].Text, text1, "the first post's text")
		Assert(t, posts[0].Author.Id, user2.Id, "first post's author")
		Assert(t, len(posts[0].Images), 0, "number of images in first post")

		Assert(t, posts[1].Text, text2, "the second post's text")
		Assert(t, posts[1].Author.Id, user2.Id, "second posts's author")
		Assert(t, len(posts[1].Images), 2, "number of images in second post")
		assertImageCreated(t, posts[1].Images[0], image1)
		assertImageCreated(t, posts[1].Images[1], image2)
		// delete them
		deletePost(t, posts[0].Id, user2)
		deletePost(t, posts[1].Id, user2)
		// assert they were deleted
		posts = getPosts(t, user2.Id, user2)
		Assert(t, len(posts), 0, "number of posts after deletion")
		assertPostFilesDeleted(t, posts[0].Id, user2.Id)
		assertPostFilesDeleted(t, posts[1].Id, user2.Id)
	})
	t.Run("liking posts", func(t *testing.T) {
		// create a post belonging to 1-st profile
		createPost(t, user1, [][]byte{}, "")
		posts := getPosts(t, user1.Id, user1)

		// like it from 2-nd profile
		toggleLike(t, posts[0].Id, user2)
		// assert it is liked
		posts = getPosts(t, user1.Id, user2)
		Assert(t, posts[0].IsLiked, true, "post is liked")

		// unlike it from 2-nd profile
		toggleLike(t, posts[0].Id, user2)
		// assert it is not liked
		posts = getPosts(t, user1.Id, user2)
		Assert(t, posts[0].IsLiked, false, "post is not liked")
	})
}

func readFixture(t testing.TB, filename string) []byte {
	t.Helper()
	return readFile(t, filepath.Join("..", "testdata", "test_avatar.jpg")) // ".." since we change the working directory to tmp_test
}

func readFile(t testing.TB, filepath string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("error while reading file %s: %v", filepath, err)
	}
	return data
}

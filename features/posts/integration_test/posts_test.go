package posts_test

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts"
	"github.com/k0marov/socnet/features/profiles"
	"github.com/k0marov/socnet/features/profiles/domain/entities"
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
	registerProfile := func(user auth.User) entities.Profile {
		fakeRegisterProfile(user)
		return entities.Profile{
			Id:       user.Id,
			Username: user.Username,
		}
	}

	// create 2 profiles
	user1 := RandomAuthUser()
	user2 := RandomAuthUser()
	profile1 := registerProfile(user1)
	profile2 := registerProfile(user2)

	t.Run("creating, reading and deleting posts", func(t *testing.T) {
		// create 1 post (without images) belonging to 2-nd profile
		text1 := "Hello, World!"
		createPost(user2.Id, [][]byte{}, text1)
		// create 1 post (with images) belonging to 2-nd profile
		text2 := "Hello, World with Images!"
		image1 := readFixture(t, "test_image.jpg")
		image2 := readFixture(t, "test_image.jpg")
		createPost(user2.Id, [][]byte{image1, image2}, text2)
		// assert they were created
		posts := getPosts(user2.Id)
		Assert(t, len(posts), 2, "number of created posts")
		//Assert(t, posts[0], )
		// delete them
	})
	t.Run("liking posts and deleting posts with images", func(t *testing.T) {
		// create a post belonging to 1-st profile
		// assert it was created
		// assert images were stored

		// like it from 2-nd profile
		// assert it is liked

		// unlike it from 2-nd profile
		// assert it is not liked

		// delete it
		// assert images were deleted

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

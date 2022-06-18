package posts_test

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/k0marov/socnet/features/posts"
	"github.com/k0marov/socnet/features/profiles"
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
	// posts
	r.Route("/posts", posts.NewPostsRouterImpl(sql, profiles.NewProfileGetterImpl(sql)))

	// helpers
	//addAuthToReq := func(req *http.Request, user core_entities.User) *http.Request {
	//	ctx := req.Context()
	//	ctx = context.WithValue(ctx, auth.UserContextKey, auth.User{Id: user.Id, Username: user.Username})
	//	return req.WithContext(ctx)
	//}

	// create 2 profiles

	// create 2 posts (without images) belonging to 2-nd profile
	// assert they were created
	// delete them

	// create a post belonging to 1-st profile
	// assert it was created
	// assert images were stored

	// like it from 2-nd profile

	// unlike it from 2-nd profile

	// delete it
	// assert images were deleted
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

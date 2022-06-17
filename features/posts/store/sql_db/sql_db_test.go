package sql_db_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/models"
	"github.com/k0marov/socnet/features/posts/store/sql_db"
	profiles_db "github.com/k0marov/socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sut, err := sql_db.NewSqlDB(db)
	AssertNoError(t, err)
	db.Close() // this will force all calls to throw errors
	t.Run("GetPosts", func(t *testing.T) {
		_, err := sut.GetPosts(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("CreatePost", func(t *testing.T) {
		_, err := sut.CreatePost(models.PostToCreate{})
		AssertSomeError(t, err)
	})
	t.Run("GetAuthor", func(t *testing.T) {
		_, err := sut.GetAuthor(RandomString())
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	t.Run("creating, reading and deleting posts", func(t *testing.T) {
		driver := OpenSqliteDB(t)

		sut, err := sql_db.NewSqlDB(driver)
		AssertNoError(t, err)
		profiles, err := profiles_db.NewSqlDB(driver)
		AssertNoError(t, err)

		createRandomPost := func(t testing.TB, author core_values.UserId) models.PostModel {
			post := models.PostToCreate{
				Author:    author,
				Text:      RandomString(),
				CreatedAt: time.Unix(time.Now().Unix(), 0),
			}
			post1Id, err := sut.CreatePost(post)
			AssertNoError(t, err)
			return models.PostModel{
				Id:        post1Id,
				Author:    author,
				Text:      post.Text,
				CreatedAt: post.CreatedAt,
				Images:    nil,
			}
		}
		assertPosts := func(t testing.TB, author core_values.UserId, posts []models.PostModel) {
			t.Helper()
			gotPosts, err := sut.GetPosts(author)
			AssertNoError(t, err)
			Assert(t, gotPosts, posts, "the stored posts")
		}
		assertAuthor := func(t testing.TB, postId values.PostId, author core_values.UserId) {
			t.Helper()
			gotAuthor, err := sut.GetAuthor(postId)
			AssertNoError(t, err)
			Assert(t, gotAuthor, author, "the stored post author")
		}

		user1 := RandomNewProfile()
		user2 := RandomNewProfile()
		profiles.CreateProfile(user1)
		profiles.CreateProfile(user2)

		wantPost1 := createRandomPost(t, user1.Id)
		assertAuthor(t, wantPost1.Id, user1.Id)
		assertPosts(t, user1.Id, []models.PostModel{wantPost1})

		wantPost1.Images = RandomPostImages()
		err = sut.AddPostImages(wantPost1.Id, wantPost1.Images)
		AssertNoError(t, err)
		assertPosts(t, user1.Id, []models.PostModel{wantPost1})

	})
}

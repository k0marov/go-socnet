package sql_db_test

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
	"time"

	"github.com/k0marov/go-socnet/features/posts/domain/models"
	"github.com/k0marov/go-socnet/features/posts/store/sql_db"
	profiles_db "github.com/k0marov/go-socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
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
	t.Run("AddPostImages", func(t *testing.T) {
		err := sut.AddPostImages(RandomString(), RandomPostImageModels())
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	createRandomPostWithTime := func(t testing.TB, sut *sql_db.SqlDB, author core_values.UserId, createdAt time.Time) models.PostModel {
		post := models.PostToCreate{
			Author:    author,
			Text:      RandomString(),
			CreatedAt: createdAt,
		}
		post1Id, err := sut.CreatePost(post)
		AssertNoError(t, err)
		return models.PostModel{
			Id:        post1Id,
			AuthorId:  author,
			Text:      post.Text,
			CreatedAt: post.CreatedAt.Unix(),
			Images:    nil,
		}
	}
	createRandomPost := func(t testing.TB, sut *sql_db.SqlDB, author core_values.UserId) models.PostModel {
		return createRandomPostWithTime(t, sut, author, RandomTime())
	}
	assertPosts := func(t testing.TB, sut *sql_db.SqlDB, author core_values.UserId, posts []models.PostModel) {
		t.Helper()
		gotPosts, err := sut.GetPosts(author)
		AssertNoError(t, err)
		Assert(t, gotPosts, posts, "the stored posts")
	}
	t.Run("creating, reading and deleting posts", func(t *testing.T) {
		driver := OpenSqliteDB(t)

		sut, err := sql_db.NewSqlDB(driver)
		AssertNoError(t, err)
		profiles, err := profiles_db.NewSqlDB(driver)
		AssertNoError(t, err)

		// create two profiles
		user1 := RandomProfileModel()
		user2 := RandomProfileModel()
		profiles.CreateProfile(user1)
		profiles.CreateProfile(user2)

		// create a post for the first profile
		wantPost1 := createRandomPost(t, sut, user1.Id)
		assertPosts(t, sut, user1.Id, []models.PostModel{wantPost1})
		// add images to that post
		wantPost1.Images = RandomPostImageModels()
		err = sut.AddPostImages(wantPost1.Id, wantPost1.Images)
		AssertNoError(t, err)
		assertPosts(t, sut, user1.Id, []models.PostModel{wantPost1})
		// create two posts for the second profile
		user2Posts := []models.PostModel{
			createRandomPost(t, sut, user2.Id),
			createRandomPost(t, sut, user2.Id),
		}
		assertPosts(t, sut, user2.Id, user2Posts)
	})
	t.Run("returning posts ordered by createdAt", func(t *testing.T) {
		driver := OpenSqliteDB(t)

		sut, err := sql_db.NewSqlDB(driver)
		AssertNoError(t, err)
		profiles, err := profiles_db.NewSqlDB(driver)
		AssertNoError(t, err)

		profile := RandomProfileModel()
		profiles.CreateProfile(profile)

		// create 3 posts
		timeInYear := func(year int) time.Time {
			return time.Date(year, 1, 1, 1, 1, 1, 0, time.UTC)
		}
		oldest := createRandomPostWithTime(t, sut, profile.Id, timeInYear(1998))
		newest := createRandomPostWithTime(t, sut, profile.Id, timeInYear(2022))
		middle := createRandomPostWithTime(t, sut, profile.Id, timeInYear(2006))
		// assert they are returned in the right order
		assertPosts(t, sut, profile.Id, []models.PostModel{newest, middle, oldest})
	})
}

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
	t.Run("AddPostImages", func(t *testing.T) {
		err := sut.AddPostImages(RandomString(), RandomPostImages())
		AssertSomeError(t, err)
	})
	t.Run("DeletePost", func(t *testing.T) {
		err := sut.DeletePost(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("IsLiked", func(t *testing.T) {
		_, err := sut.IsLiked(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("LikePost", func(t *testing.T) {
		err := sut.LikePost(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("UnlikePost", func(t *testing.T) {
		err := sut.UnlikePost(RandomString(), RandomString())
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
			Author:    author,
			Text:      post.Text,
			CreatedAt: post.CreatedAt,
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

		assertAuthor := func(t testing.TB, postId values.PostId, author core_values.UserId) {
			t.Helper()
			gotAuthor, err := sut.GetAuthor(postId)
			AssertNoError(t, err)
			Assert(t, gotAuthor, author, "the stored post author")
		}
		// create two profiles
		user1 := RandomNewProfile()
		user2 := RandomNewProfile()
		profiles.CreateProfile(user1)
		profiles.CreateProfile(user2)

		// create a post for the first profile
		wantPost1 := createRandomPost(t, sut, user1.Id)
		assertAuthor(t, wantPost1.Id, user1.Id)
		assertPosts(t, sut, user1.Id, []models.PostModel{wantPost1})
		// add images to that post
		wantPost1.Images = RandomPostImages()
		err = sut.AddPostImages(wantPost1.Id, wantPost1.Images)
		AssertNoError(t, err)
		assertPosts(t, sut, user1.Id, []models.PostModel{wantPost1})
		// create two posts for the second profile
		user2Posts := []models.PostModel{
			createRandomPost(t, sut, user2.Id),
			createRandomPost(t, sut, user2.Id),
		}
		assertAuthor(t, user2Posts[0].Id, user2.Id)
		assertAuthor(t, user2Posts[1].Id, user2.Id)
		assertPosts(t, sut, user2.Id, user2Posts)

		// delete the second post
		err = sut.DeletePost(user2Posts[1].Id)
		AssertNoError(t, err)
		// assert it was deleted
		assertPosts(t, sut, user2.Id, user2Posts[:1])
	})
	t.Run("liking posts", func(t *testing.T) {
		driver := OpenSqliteDB(t)

		sut, err := sql_db.NewSqlDB(driver)
		AssertNoError(t, err)
		profiles, err := profiles_db.NewSqlDB(driver)
		AssertNoError(t, err)

		assertIsLiked := func(t testing.TB, post values.PostId, byUser core_values.UserId, isLiked bool) {
			t.Helper()
			gotIsLiked, err := sut.IsLiked(post, byUser)
			AssertNoError(t, err)
			Assert(t, gotIsLiked, isLiked, "the returned isLiked")
		}

		// create two profiles
		user1 := RandomNewProfile()
		user2 := RandomNewProfile()
		profiles.CreateProfile(user1)
		profiles.CreateProfile(user2)

		// create a random post belonging to user1
		post := createRandomPost(t, sut, user1.Id)
		post2 := createRandomPost(t, sut, user2.Id)
		// it shouldn't be liked by any of the users
		assertIsLiked(t, post.Id, user1.Id, false)
		assertIsLiked(t, post.Id, user2.Id, false)
		// like it from user2
		err = sut.LikePost(post.Id, user2.Id)
		AssertNoError(t, err)
		// now it should be liked from user2 and not liked from user1
		assertIsLiked(t, post.Id, user1.Id, false)
		assertIsLiked(t, post.Id, user2.Id, true)
		// like it from user1
		err = sut.LikePost(post.Id, user1.Id)
		AssertNoError(t, err)
		// now it should be liked by both users
		assertIsLiked(t, post.Id, user1.Id, true)
		assertIsLiked(t, post.Id, user2.Id, true)

		// unlike it from user2
		err = sut.UnlikePost(post.Id, user2.Id)
		AssertNoError(t, err)
		// now it should be liked only by user1
		assertIsLiked(t, post.Id, user1.Id, true)
		assertIsLiked(t, post.Id, user2.Id, false)

		// post2 should not be affected
		assertIsLiked(t, post2.Id, user1.Id, false)
		assertIsLiked(t, post2.Id, user2.Id, false)
	})
	t.Run("returning posts ordered by createdAt", func(t *testing.T) {
		driver := OpenSqliteDB(t)

		sut, err := sql_db.NewSqlDB(driver)
		AssertNoError(t, err)
		profiles, err := profiles_db.NewSqlDB(driver)
		AssertNoError(t, err)

		profile := RandomNewProfile()
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

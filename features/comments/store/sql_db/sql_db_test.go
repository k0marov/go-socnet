package sql_db_test

import (
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/domain/values"
	"github.com/k0marov/socnet/features/comments/store/models"
	"github.com/k0marov/socnet/features/comments/store/sql_db"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/post_models"
	posts_db "github.com/k0marov/socnet/features/posts/store/sql_db"
	profiles_db "github.com/k0marov/socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB, err := sql_db.NewSqlDB(db)
	AssertNoError(t, err)
	db.Close() // this will make all calls to db throw
	t.Run("IsLiked", func(t *testing.T) {
		_, err := sqlDB.IsLiked(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Like", func(t *testing.T) {
		err := sqlDB.Like(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Unlike", func(t *testing.T) {
		err := sqlDB.Unlike(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("GetComments", func(t *testing.T) {
		_, err := sqlDB.GetComments(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Create", func(t *testing.T) {
		_, err := sqlDB.Create(RandomNewComment(), RandomTime())
		AssertSomeError(t, err)
	})
	t.Run("GetAuthor", func(t *testing.T) {
		_, err := sqlDB.GetAuthor(RandomId())
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	createComment := func(t testing.TB, db *sql_db.SqlDB, post post_values.PostId, author core_values.UserId, yearCreatedAt int) models.CommentModel {
		t.Helper()
		newComment := values.NewCommentValue{
			Author: author,
			Post:   post,
			Text:   RandomString(),
		}
		createdAt := time.Date(yearCreatedAt, 0, 0, 0, 0, 0, 0, time.UTC)
		id, err := db.Create(newComment, createdAt)
		AssertNoError(t, err)
		return models.CommentModel{
			Id:        id,
			Author:    author,
			Text:      newComment.Text,
			CreatedAt: createdAt,
			Likes:     0,
		}
	}
	getComments := func(t testing.TB, db *sql_db.SqlDB, post post_values.PostId) []models.CommentModel {
		t.Helper()
		comments, err := db.GetComments(post)
		AssertNoError(t, err)
		return comments
	}

	t.Run("creating and reading comments", func(t *testing.T) {
		assertAuthor := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, author core_values.UserId) {
			t.Helper()
			gotAuthor, err := db.GetAuthor(comment)
			AssertNoError(t, err)
			Assert(t, gotAuthor, author, "author")
		}

		db := OpenSqliteDB(t)
		sqlDB, err := sql_db.NewSqlDB(db)
		AssertNoError(t, err)
		profilesDb, _ := profiles_db.NewSqlDB(db)
		postsDb, _ := posts_db.NewSqlDB(db)

		// create an author profile
		author := RandomNewProfile()
		profilesDb.CreateProfile(author)

		// create a post
		postId, _ := postsDb.CreatePost(post_models.PostToCreate{
			Author:    author.Id,
			Text:      RandomString(),
			CreatedAt: time.Now(),
		})

		// create a commenter profile
		commenter := RandomNewProfile()
		profilesDb.CreateProfile(commenter)

		// create the first comment
		firstComment := createComment(t, sqlDB, postId, commenter.Id, 2020)

		// assert it was created
		comments := getComments(t, sqlDB, postId)
		AssertFatal(t, len(comments), 1, "number of post comments")
		Assert(t, comments[0], firstComment, "the created comment")
		assertAuthor(t, sqlDB, comments[0].Id, commenter.Id)

		// create the second comment
		secondComment := createComment(t, sqlDB, postId, commenter.Id, 2022)

		// assert it was created (and comments are returned ordered by createdAt)
		comments = getComments(t, sqlDB, postId)
		AssertFatal(t, len(comments), 2, "number of post comments")
		assertAuthor(t, sqlDB, comments[1].Id, commenter.Id)
		Assert(t, comments[0], secondComment, "the second created comment")
		Assert(t, comments[1], firstComment, "the first created comment")
	})
	t.Run("liking comments", func(t *testing.T) {
		assertLikedValue := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, liker core_values.UserId, value bool) {
			t.Helper()
			isLiked, err := db.IsLiked(comment, liker)
			AssertNoError(t, err)
			Assert(t, isLiked, value, "isLiked value")
		}
		assertNotLiked := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, liker core_values.UserId) {
			t.Helper()
			assertLikedValue(t, db, comment, liker, false)
		}
		assertLiked := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, liker core_values.UserId) {
			t.Helper()
			assertLikedValue(t, db, comment, liker, true)
		}
		unlike := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, unliker core_values.UserId) {
			t.Helper()
			err := db.Unlike(comment, unliker)
			AssertNoError(t, err)
		}
		like := func(t testing.TB, db *sql_db.SqlDB, comment values.CommentId, liker core_values.UserId) {
			t.Helper()
			err := db.Like(comment, liker)
			AssertNoError(t, err)
		}

		db := OpenSqliteDB(t)
		sqlDB, err := sql_db.NewSqlDB(db)
		AssertNoError(t, err)
		profilesDb, _ := profiles_db.NewSqlDB(db)
		postsDb, _ := posts_db.NewSqlDB(db)

		// create an author profile
		author := RandomNewProfile()
		profilesDb.CreateProfile(author)

		// create a post
		postId, _ := postsDb.CreatePost(post_models.PostToCreate{
			Author:    author.Id,
			Text:      RandomString(),
			CreatedAt: time.Now(),
		})

		// create a comment
		commenter := RandomNewProfile()
		profilesDb.CreateProfile(commenter)
		comment := createComment(t, sqlDB, postId, commenter.Id, 2020)

		// create another profile
		liker := RandomNewProfile()
		profilesDb.CreateProfile(liker)
		// assert the comment is not liked from this profile
		assertNotLiked(t, sqlDB, comment.Id, liker.Id)
		// like it
		like(t, sqlDB, comment.Id, liker.Id)
		// assert it was liked
		assertLiked(t, sqlDB, comment.Id, liker.Id)
		// unlike it
		unlike(t, sqlDB, comment.Id, liker.Id)
		// assert it is not liked
		assertNotLiked(t, sqlDB, comment.Id, liker.Id)

		// liking a comment from many profiles
		const count = 100
		for i := 0; i < count; i++ {
			profile := RandomNewProfile()
			profilesDb.CreateProfile(profile)
			like(t, sqlDB, comment.Id, profile.Id)
			comments := getComments(t, sqlDB, postId)
			AssertFatal(t, len(comments), 1, "number of comments")
			Assert(t, comments[0].Likes, i+1, "number of likes")
		}
	})
}

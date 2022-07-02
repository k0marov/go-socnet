package sql_db_test

import (
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"
	"time"

	"github.com/k0marov/go-socnet/features/comments/domain/models"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	"github.com/k0marov/go-socnet/features/comments/store/sql_db"
	post_models "github.com/k0marov/go-socnet/features/posts/domain/models"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
	posts_db "github.com/k0marov/go-socnet/features/posts/store/sql_db"
	profiles_db "github.com/k0marov/go-socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB, err := sql_db.NewSqlDB(db)
	AssertNoError(t, err)
	db.Close() // this will make all calls to db throw
	t.Run("GetComments", func(t *testing.T) {
		_, err := sqlDB.GetComments(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Create", func(t *testing.T) {
		_, err := sqlDB.Create(RandomNewComment(), RandomTime())
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
			AuthorId:  author,
			Text:      newComment.Text,
			CreatedAt: createdAt,
		}
	}
	getComments := func(t testing.TB, db *sql_db.SqlDB, post post_values.PostId) []models.CommentModel {
		t.Helper()
		comments, err := db.GetComments(post)
		AssertNoError(t, err)
		return comments
	}

	t.Run("creating and reading comments", func(t *testing.T) {
		db := OpenSqliteDB(t)
		sqlDB, err := sql_db.NewSqlDB(db)
		AssertNoError(t, err)
		profilesDb, _ := profiles_db.NewSqlDB(db)
		postsDb, _ := posts_db.NewSqlDB(db)

		// create an author profile
		author := RandomProfileModel()
		profilesDb.CreateProfile(author)

		// create a post
		postId, _ := postsDb.CreatePost(post_models.PostToCreate{
			Author:    author.Id,
			Text:      RandomString(),
			CreatedAt: time.Now(),
		})

		// create a commenter profile
		commenter := RandomProfileModel()
		profilesDb.CreateProfile(commenter)

		// create the first comment
		firstComment := createComment(t, sqlDB, postId, commenter.Id, 2020)

		// assert it was created
		comments := getComments(t, sqlDB, postId)
		AssertFatal(t, len(comments), 1, "number of post comments")
		Assert(t, comments[0], firstComment, "the created comment")

		// create the second comment
		secondComment := createComment(t, sqlDB, postId, commenter.Id, 2022)

		// assert it was created (and comments are returned ordered by createdAt)
		comments = getComments(t, sqlDB, postId)
		AssertFatal(t, len(comments), 2, "number of post comments")
		Assert(t, comments[0], secondComment, "the second created comment")
		Assert(t, comments[1], firstComment, "the first created comment")
	})
}

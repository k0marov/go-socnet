package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/k0marov/go-socnet/core/general/core_values"
	helpers "github.com/k0marov/go-socnet/core/helpers/http_test_helpers"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/k0marov/go-socnet/features/comments"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/comments/delivery/http/responses"
	"github.com/k0marov/go-socnet/features/comments/domain/values"
	post_models "github.com/k0marov/go-socnet/features/posts/domain/models"
	post_values "github.com/k0marov/go-socnet/features/posts/domain/values"
	posts_db "github.com/k0marov/go-socnet/features/posts/store/sql_db"
	"github.com/k0marov/go-socnet/features/profiles"
	profile_responses "github.com/k0marov/go-socnet/features/profiles/delivery/http/responses"
	auth "github.com/k0marov/golang-auth"
	_ "github.com/mattn/go-sqlite3"
)

func TestComments(t *testing.T) {
	// db
	sql, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error while opening in-memory database: %v", err)
	}
	r := chi.NewRouter()
	// profiles
	fakeRegisterProfile := profiles.NewRegisterCallback(sql)
	getProfile := profiles.NewProfileGetterImpl(sql)
	// posts
	postsDB, _ := posts_db.NewSqlDB(sql)
	createPost := func(author core_values.UserId) post_values.PostId {
		id, _ := postsDB.CreatePost(post_models.PostToCreate{
			Author:    author,
			Text:      RandomString(),
			CreatedAt: RandomTime(),
		})
		return id
	}
	// comments
	r.Route("/comments", comments.NewCommentsRouterImpl(sql, getProfile))

	assertComments := func(t testing.TB, got, want []responses.CommentResponse) {
		t.Helper()
		AssertFatal(t, len(got), len(want), "number of returned comments")
		for i, comment := range got {
			Assert(t, comment, want[i], "comment")
		}
	}
	addComment := func(t testing.TB, post post_values.PostId, caller auth.User) responses.CommentResponse {
		t.Helper()

		newComment := handlers.NewCommentRequest{Text: RandomString()}
		body := bytes.NewBuffer(nil)
		json.NewEncoder(body).Encode(newComment)
		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodPost, "/comments/?post_id="+post, body), caller)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		var returnedComment responses.CommentResponse
		json.NewDecoder(response.Body).Decode(&returnedComment)

		wantAuthor, _ := getProfile(caller.Id, caller.Id)
		wantComment := responses.CommentResponse{
			Id:        returnedComment.Id,
			Author:    profile_responses.NewProfileResponse(wantAuthor),
			Text:      newComment.Text,
			CreatedAt: time.Now().Unix(),
			Likes:     0,
			IsLiked:   false,
			IsMine:    true,
		}
		assertComments(t, []responses.CommentResponse{returnedComment}, []responses.CommentResponse{wantComment})

		return returnedComment
	}
	getComments := func(t testing.TB, post post_values.PostId, caller auth.User) []responses.CommentResponse {
		t.Helper()

		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodGet, "/comments/?post_id="+post, nil), caller)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		var commentsResponse responses.CommentsResponse
		json.NewDecoder(response.Body).Decode(&commentsResponse)

		return commentsResponse.Comments
	}
	toggleLike := func(t testing.TB, comment values.CommentId, caller auth.User) {
		t.Helper()
		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodPost, "/comments/"+comment+"/toggle-like", nil), caller)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK)
	}
	deleteComment := func(t testing.TB, comment values.CommentId, caller auth.User) {
		t.Helper()
		request := helpers.AddAuthDataToRequest(httptest.NewRequest(http.MethodDelete, "/comments/"+comment, nil), caller)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK)
	}

	t.Run("creating, reading and deleting comments", func(t *testing.T) {

		// create 2 profiles
		user1 := RandomAuthUser()
		user2 := RandomAuthUser()
		fakeRegisterProfile(user1)
		fakeRegisterProfile(user2)

		// create post belonging to 1-st profile
		post := createPost(user1.Id)
		// add a comment to this post from 2-nd profile
		comment1 := addComment(t, post, user2)
		// assert it was created
		comments := getComments(t, post, user2)
		assertComments(t, comments, []responses.CommentResponse{comment1})

		// wait so that createdAt difference is at least 1 second
		time.Sleep(time.Second)

		// create another comment from 2-nd profile
		comment2 := addComment(t, post, user2)
		// assert it was created and comments are returned in the right order (newest first)
		comments = getComments(t, post, user2)
		assertComments(t, comments, []responses.CommentResponse{comment2, comment1})

		t.Run("liking comments", func(t *testing.T) {
			// like the newest comment from 1-st profile
			toggleLike(t, comment2.Id, user1)
			// assert it is liked
			comments = getComments(t, post, user1)
			Assert(t, comments[0].IsLiked, true, "isLiked")
			// unlike it
			toggleLike(t, comment2.Id, user1)
			// assert it is not liked
			comments = getComments(t, post, user1)
			Assert(t, comments[0].IsLiked, false, "isLiked")
		})
		// delete the second comment
		deleteComment(t, comment2.Id, user2)
		// assert it was deleted
		comments = getComments(t, post, user2)
		assertComments(t, comments, []responses.CommentResponse{comment1})
		// delete the first comment
		deleteComment(t, comment1.Id, user2)
		// assert it was deleted
		assertComments(t, getComments(t, post, user2), []responses.CommentResponse{})
	})
}

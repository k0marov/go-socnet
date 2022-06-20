package test_helpers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/ref"
	comment_entities "github.com/k0marov/socnet/features/comments/domain/entities"
	post_values "github.com/k0marov/socnet/features/posts/domain/values"
	"github.com/k0marov/socnet/features/posts/store/models"
	"math"
	random "math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	profile_entities "github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"

	post_entities "github.com/k0marov/socnet/features/posts/domain/entities"

	"github.com/k0marov/socnet/core/client_errors"
	core_entities "github.com/k0marov/socnet/core/core_entities"

	auth "github.com/k0marov/golang-auth"
)

var rand = random.New(random.NewSource(time.Now().UnixNano()))

func AssertStatusCode(t testing.TB, got *httptest.ResponseRecorder, want int) {
	t.Helper()
	Assert(t, got.Result().StatusCode, want, "response status code")
}

func Assert[T any](t testing.TB, got, want T, description string) bool {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s is not right: got '%v', want '%v'", description, got, want)
		return false
	}
	return true
}

func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("expected no error but got %v", got)
	}
}
func AssertError(t testing.TB, got error, want error) {
	t.Helper()
	if got != want {
		t.Errorf("expected error %v, but got %v", want, got)
	}
}
func AssertSomeError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Error("expected an error, but got nil")
	}
}

func AssertFatal[T comparable](t testing.TB, got, want T, description string) {
	t.Helper()
	if !Assert(t, got, want, description) {
		t.Fatal()
	}
}

func AssertClientError(t testing.TB, response *httptest.ResponseRecorder, err client_errors.ClientError) {
	t.Helper()
	var got client_errors.ClientError
	json.NewDecoder(response.Body).Decode(&got)

	AssertJSON(t, response)
	Assert(t, got, err, "error response")
	Assert(t, response.Code, err.HTTPCode, "status code")
}

func AssertJSON(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	Assert(t, response.Result().Header.Get("contentType"), "application/json", "response content type")
}

func AssertJSONData[T any](t testing.TB, response *httptest.ResponseRecorder, wantData T) {
	t.Helper()
	AssertStatusCode(t, response, http.StatusOK)
	AssertJSON(t, response)
	var gotData T
	json.NewDecoder(response.Body).Decode(&gotData)
	Assert(t, gotData, wantData, "json encoded data")
}

func RandomError() error {
	return errors.New(RandomString())
}

func RandomUser() core_entities.User {
	return core_entities.User{
		Id:       strconv.Itoa(RandomInt()),
		Username: RandomString(),
	}
}

func RandomAuthUser() auth.User {
	return auth.User{
		Id:       strconv.Itoa(RandomInt()),
		Username: RandomString(),
	}
}

func RandomProfile() profile_entities.Profile {
	return profile_entities.Profile{
		Id:         strconv.Itoa(RandomInt()),
		Username:   RandomString(),
		About:      RandomString(),
		AvatarPath: RandomString(),
		Follows:    RandomInt(),
		Followers:  RandomInt(),
	}
}

func RandomContextedProfile() profile_entities.ContextedProfile {
	return profile_entities.ContextedProfile{
		Profile:            RandomProfile(),
		IsFollowedByCaller: RandomBool(),
	}
}

func RandomNewProfile() values.NewProfile {
	return values.NewProfile{
		Id:         strconv.Itoa(RandomInt()),
		Username:   RandomString(),
		About:      RandomString(),
		AvatarPath: RandomString(),
	}
}

func RandomFiles() []core_values.FileData {
	return []core_values.FileData{RandomFileData(), RandomFileData(), RandomFileData()}
}
func RandomUrls() []core_values.FileURL {
	return []core_values.FileURL{RandomString(), RandomString(), RandomString()}
}
func RandomPostImages() []post_values.PostImage {
	return []post_values.PostImage{
		{URL: RandomString(), Index: 1},
		{URL: RandomString(), Index: 2},
		{URL: RandomString(), Index: 3},
	}
}

func RandomContextedPost() post_entities.ContextedPost {
	return post_entities.ContextedPost{
		Id:        strconv.Itoa(RandomInt()),
		Author:    RandomContextedProfile(),
		Text:      RandomString(),
		Images:    RandomPostImages(),
		CreatedAt: RandomTime(),
		IsMine:    RandomBool(),
		IsLiked:   RandomBool(),
	}
}
func RandomNewPostData() post_values.NewPostData {
	return post_values.NewPostData{
		Text:   RandomString(),
		Author: RandomString(),
		Images: []post_values.PostImageFile{{RandomFileData(), 1}, {RandomFileData(), 2}},
	}
}

func RandomTime() time.Time {
	return time.Date(2022, 6, 17, 16, 53, 42, 0, time.UTC)
}

func RandomPostModel() models.PostModel {
	return models.PostModel{
		Id:        RandomString(),
		Author:    RandomString(),
		Text:      RandomString(),
		CreatedAt: RandomTime(),
		Images:    RandomPostImages(),
	}
}

func OpenSqliteDB(t testing.TB) *sql.DB {
	t.Helper()
	sql, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error while opening in-memory database: %v", err)
	}
	return sql
}

func TimeAlmostEqual(t, want time.Time) bool {
	return math.Abs(t.Sub(want).Minutes()) < 1
}

func TimeAlmostNow(t time.Time) bool {
	return TimeAlmostEqual(t, time.Now())
}

func RandomFileData() core_values.FileData {
	data := []byte(RandomString())
	ref, _ := ref.NewRef(&data)
	return ref
}

func RandomClientError() client_errors.ClientError {
	return client_errors.ClientError{
		DetailCode:     RandomString(),
		ReadableDetail: RandomString(),
		HTTPCode:       RandomInt() + 400,
	}
}

func RandomComment() comment_entities.Comment {
	return comment_entities.Comment{}
}
func RandomComments() []comment_entities.Comment {
	return []comment_entities.Comment{RandomComment(), RandomComment(), RandomComment()}
}

func RandomBool() bool {
	return rand.Float32() > 0.5
}

func RandomInt() int {
	return rand.Intn(100)
}

func RandomString() string {
	str := ""
	for i := 0; i < 2; i++ {
		str += words[rand.Intn(len(words))] + "_"
	}
	return str
}

var words = []string{"the", "be", "to", "of", "and", "a", "in", "that", "have", "I", "it", "for", "not", "on", "with", "he", "as", "you", "do", "at", "this", "but", "his", "by", "from", "they", "we", "say", "her", "she", "or", "an", "will", "my", "one", "all", "would", "there", "their", "what", "so", "up", "out", "if", "about", "who", "get", "which", "go", "me", "when", "make", "can", "like", "time", "no", "just", "him", "know", "take", "people", "into", "year", "your", "good", "some", "could", "them", "see", "other", "than", "then", "now", "look", "only", "come", "its", "over", "think", "also", "back", "after", "use", "two", "how", "our", "work", "first", "well", "way", "even", "new", "want", "because", "any", "these", "give", "day", "most", "us"}

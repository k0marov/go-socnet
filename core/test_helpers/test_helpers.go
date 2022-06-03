package test_helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"os"
	"profiles/domain/entities"
	"reflect"
	"testing"

	"core/client_errors"
	core_entities "core/entities"

	auth "github.com/k0marov/golang-auth"
)

func CreateTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

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
		t.Errorf("expected no error but got %v", got)
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

func AssertNotNil[T comparable](t testing.TB, got T, description string) {
	t.Helper()
	var nilT T
	if got == nilT {
		t.Errorf("expected %s to be non nil, but got nil", description)
	}
}

func AssertHTTPError(t testing.TB, response *httptest.ResponseRecorder, err client_errors.ClientError, statusCode int) {
	t.Helper()
	var got client_errors.ClientError
	json.NewDecoder(response.Body).Decode(&got)

	AssertJSON(t, response)
	Assert(t, got, err, "error response")
	Assert(t, response.Code, statusCode, "status code")
}

func AssertJSON(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()
	Assert(t, response.Result().Header.Get("contentType"), "application/json", "response content type")
}

func AssertJSONData[T any](t testing.TB, response *httptest.ResponseRecorder, wantData T) {
	t.Helper()
	AssertJSON(t, response)
	var gotData T
	json.NewDecoder(response.Body).Decode(&gotData)
	Assert(t, gotData, wantData, "json encoded data")
}

func AssertUniqueCount[T comparable](t testing.TB, slice []T, want int) {
	t.Helper()
	unique := []T{}
	for _, val := range slice {
		if !CheckInSlice(val, unique) {
			unique = append(unique, val)
		}
	}
	Assert(t, len(unique), want, "number of unique elements")
}

func CheckInSlice[T comparable](elem T, slice []T) bool {
	for _, sliceElem := range slice {
		if sliceElem == elem {
			return true
		}
	}
	return false
}

func RandomError() error {
	return errors.New(RandomString())
}

func RandomUser() core_entities.User {
	return core_entities.User{
		Id:       RandomString(),
		Username: RandomString(),
	}
}

func RandomAuthUser() auth.User {
	return auth.User{
		Id:       RandomString(),
		Username: RandomString(),
	}
}

func RandomDetailedProfile() entities.DetailedProfile {
	return entities.DetailedProfile{
		Profile: RandomProfile(),
	}
}

func RandomProfile() entities.Profile {
	return entities.Profile{
		Id:       RandomString(),
		Username: RandomString(),
	}
}

func RandomString() string {
	str := ""
	for i := 0; i < 2; i++ {
		str += words[rand.Intn(len(words))] + "_"
	}
	return str
}

var words = []string{"the", "be", "to", "of", "and", "a", "in", "that", "have", "I", "it", "for", "not", "on", "with", "he", "as", "you", "do", "at", "this", "but", "his", "by", "from", "they", "we", "say", "her", "she", "or", "an", "will", "my", "one", "all", "would", "there", "their", "what", "so", "up", "out", "if", "about", "who", "get", "which", "go", "me", "when", "make", "can", "like", "time", "no", "just", "him", "know", "take", "people", "into", "year", "your", "good", "some", "could", "them", "see", "other", "than", "then", "now", "look", "only", "come", "its", "over", "think", "also", "back", "after", "use", "two", "how", "our", "work", "first", "well", "way", "even", "new", "want", "because", "any", "these", "give", "day", "most", "us"}

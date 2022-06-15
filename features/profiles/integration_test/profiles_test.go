package profiles_test

import (
	"bytes"
	"context"
	"core/core_values"
	core_entities "core/entities"
	. "core/test_helpers"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"profiles"
	"profiles/delivery/http/handlers"
	"profiles/domain/entities"
	"profiles/domain/values"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

func TestProfiles(t *testing.T) {
	// working directory
	os.Mkdir("tmp_test", 0777)
	os.Chdir("tmp_test")
	defer func() {
		os.Chdir("..")
		os.RemoveAll("tmp_test")
	}()

	// profiles setup
	sql, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error while opening in-memory database: %v", err)
	}

	r := chi.NewRouter()
	r.Route("/profiles", profiles.NewProfilesRouterImpl(sql))

	// fake auth setup
	fakeRegisterRequest := func(newUser core_entities.User) { // mock registering a new user
		callback := profiles.NewRegisterCallback(sql)
		callback(auth.User{Id: newUser.Id, Username: newUser.Username})
	}

	// helpers
	addAuthToReq := func(req *http.Request, user core_entities.User) *http.Request {
		ctx := req.Context()
		ctx = context.WithValue(ctx, auth.UserContextKey, auth.User{Id: user.Id, Username: user.Username})
		return req.WithContext(ctx)
	}
	checkProfileFromServer := func(t testing.TB, wantProfile entities.DetailedProfile) {
		t.Helper()
		request := addAuthToReq(httptest.NewRequest(http.MethodGet, "/profiles/"+wantProfile.Id, nil), RandomUser())
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertJSONData(t, response, wantProfile.Profile)

		request = addAuthToReq(httptest.NewRequest(http.MethodGet, "/profiles/me", nil), core_entities.User{Id: wantProfile.Id, Username: wantProfile.Username})
		response = httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertJSONData(t, response, wantProfile)
	}
	t.Run("creating, reading and updating", func(t *testing.T) {
		// register a couple users
		user1 := RandomUser()
		user2 := RandomUser()
		fakeRegisterRequest(user1)
		fakeRegisterRequest(user2)

		// assert that users are now accessible from GET handler
		profile1 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         user1.Id,
				Username:   user1.Username,
				About:      "",
				AvatarPath: "",
				Follows:    0,
				Followers:  0,
			},
			FollowsProfiles: []entities.Profile{},
		}
		profile2 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         user2.Id,
				Username:   user2.Username,
				About:      "",
				AvatarPath: "",
				Follows:    0,
				Followers:  0,
			},
			FollowsProfiles: []entities.Profile{},
		}
		checkProfileFromServer(t, profile1)
		checkProfileFromServer(t, profile2)

		// update avatar for first user
		wantAvatarPathStr := filepath.Join("static", "user_"+user1.Id, "avatar")
		wantAvatarPath := values.AvatarPath{Path: wantAvatarPathStr}
		avatar := readFixture(t, "test_avatar.jpg")

		body, contentType := createMultipartBody(avatar)
		request := addAuthToReq(httptest.NewRequest(http.MethodPut, "/profiles/me/avatar", body), user1)
		request.Header.Add("Content-Type", contentType)
		response := httptest.NewRecorder()

		r.ServeHTTP(response, request)
		AssertJSONData(t, response, wantAvatarPath)

		// assert that it was updated
		wantUpdatedProfile1 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         user1.Id,
				Username:   user1.Username,
				About:      "",
				AvatarPath: wantAvatarPathStr,
				Follows:    0,
				Followers:  0,
			},
			FollowsProfiles: []entities.Profile{},
		}
		checkProfileFromServer(t, wantUpdatedProfile1)

		// assert avatar was stored
		Assert(t, readFile(t, wantAvatarPathStr), avatar, "the stored avatar file")

		// update profile for second user
		upd := values.ProfileUpdateData{About: RandomString()}
		reqBody := bytes.NewBuffer(nil)
		json.NewEncoder(reqBody).Encode(upd)
		request = addAuthToReq(httptest.NewRequest(http.MethodPut, "/profiles/me", reqBody), user2)
		response = httptest.NewRecorder()

		r.ServeHTTP(response, request)
		wantUpdatedProfile2 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         user2.Id,
				Username:   user2.Username,
				About:      upd.About,
				AvatarPath: "",
				Follows:    0,
				Followers:  0,
			},
			FollowsProfiles: []entities.Profile{},
		}
		AssertJSONData(t, response, wantUpdatedProfile2)

		// check it was stored
		checkProfileFromServer(t, wantUpdatedProfile2)

	})
	t.Run("following", func(t *testing.T) {
		checkFollows := func(t testing.TB, id core_values.UserId, wantFollows []entities.Profile) {
			t.Helper()
			request := addAuthToReq(httptest.NewRequest(http.MethodGet, "/profiles/"+id+"/follows", nil), RandomUser())
			response := httptest.NewRecorder()
			r.ServeHTTP(response, request)
			AssertJSONData(t, response, handlers.FollowsResponse{Profiles: wantFollows})
		}

		// create 2 users
		user1 := RandomUser()
		user2 := RandomUser()
		fakeRegisterRequest(user1)
		fakeRegisterRequest(user2)

		// follow profile1 from profile2
		request := addAuthToReq(httptest.NewRequest(http.MethodPost, "/profiles/"+user1.Id+"/toggle-follow", nil), user2)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK)

		// assert it was followed
		wantProfile1 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:        user1.Id,
				Username:  user1.Username,
				Follows:   0,
				Followers: 1,
			},
			FollowsProfiles: []entities.Profile{},
		}
		checkProfileFromServer(t, wantProfile1)

		wantFollows := []entities.Profile{wantProfile1.Profile}
		checkFollows(t, user2.Id, wantFollows)
		wantProfile2 := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:        user2.Id,
				Username:  user2.Username,
				Follows:   1,
				Followers: 0,
			},
			FollowsProfiles: wantFollows,
		}
		checkProfileFromServer(t, wantProfile2)

		// unfollow profile1 from profile2
		request = addAuthToReq(httptest.NewRequest(http.MethodPost, "/profiles/"+user1.Id+"/toggle-follow", nil), user2)
		response = httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertStatusCode(t, response, http.StatusOK)

		// assert it is now not followed
		wantProfile1.Profile.Followers = 0
		checkProfileFromServer(t, wantProfile1)
		wantFollows = []entities.Profile{}
		checkFollows(t, user2.Id, wantFollows)
		wantProfile2.Profile.Follows = 0
		wantProfile2.FollowsProfiles = wantFollows
		checkProfileFromServer(t, wantProfile2)
	})

}

func readFixture(t testing.TB, filename string) []byte {
	t.Helper()
	return readFile(t, filepath.Join("..", "testdata", "test_avatar.jpg")) // ".." since we change the working directory to tmp_test
}

func readFile(t testing.TB, filepath string) []byte {
	data, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("error while reading file %s: %v", filepath, err)
	}
	return data
}

func createMultipartBody(data []byte) (io.Reader, string) {
	body := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	fw, _ := writer.CreateFormFile("avatar", RandomString())
	fw.Write(data)
	return body, writer.FormDataContentType()
}

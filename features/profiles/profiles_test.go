package profiles_test

import (
	"context"
	core_entities "core/entities"
	. "core/test_helpers"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"profiles"
	"profiles/domain/entities"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi/v5"
	auth "github.com/k0marov/golang-auth"
)

func TestProfiles(t *testing.T) {
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
		},
	}
	profile2 := entities.DetailedProfile{
		Profile: entities.Profile{
			Id:         user2.Id,
			Username:   user2.Username,
			About:      "",
			AvatarPath: "",
		},
	}
	checkGetMeRequestForUser := func(fromUser core_entities.User, wantProfile entities.DetailedProfile) {
		request := addAuthToReq(httptest.NewRequest(http.MethodGet, "/profiles/me", nil), fromUser)
		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)
		AssertJSONData(t, response, wantProfile)
	}
	checkGetMeRequestForUser(user1, profile1)
	checkGetMeRequestForUser(user2, profile2)

	// TODO: update avatar for first user

}

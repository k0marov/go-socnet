package handlers_test

import (
	"context"
	"core/client_errors"
	core_entities "core/entities"
	. "core/test_helpers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/http/handlers"
	"profiles/domain/entities"
	"profiles/domain/values"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetMeHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	createRequestWithAuth := func() *http.Request {
		return addAuthDataToRequest(createRequest(nil), authUser)
	}
	baseTest401(t, handlers.NewGetMeHandler(nil))
	t.Run("should return 200 and a profile if authentication details are provided via context", func(t *testing.T) {
		wantedProfile := RandomDetailedProfile()
		getter := func(gotUser core_entities.User) (entities.DetailedProfile, error) {
			if gotUser == user {
				return wantedProfile, nil
			}
			panic(fmt.Sprintf("called with user=%v", gotUser))
		}

		response := httptest.NewRecorder()
		handlers.NewGetMeHandler(getter).ServeHTTP(response, createRequestWithAuth())

		AssertJSONData(t, response, wantedProfile)
	})
	baseTestServiceErrorHandling(t, func(wantErr error, response *httptest.ResponseRecorder) {
		getter := func(core_entities.User) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, wantErr
		}
		handlers.NewGetMeHandler(getter).ServeHTTP(response, createRequestWithAuth())
	})
}
func createRequestWithId(userId string) *http.Request {
	request := createRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", userId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return request
}

func TestGetByIdHandler(t *testing.T) {
	t.Run("should return 200 and a profile if profile with given id exists", func(t *testing.T) {
		randomId := RandomString()
		randomProfile := RandomProfile()
		profileGetter := func(userId values.UserId) (entities.Profile, error) {
			if userId == randomId {
				return randomProfile, nil
			}
			panic("called with unexpected arguments")
		}

		request := createRequestWithId(randomId)
		response := httptest.NewRecorder()

		handlers.NewGetByIdHandler(profileGetter).ServeHTTP(response, request)
		AssertJSONData(t, response, randomProfile)
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		handlers.NewGetByIdHandler(nil).ServeHTTP(response, createRequest(nil)) // getter is nil, since it shouldn't be called
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	baseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(userId values.UserId) (entities.Profile, error) {
			return entities.Profile{}, err
		}
		handlers.NewGetByIdHandler(getter).ServeHTTP(rr, createRequestWithId("42"))
	})
}

func TestFollowsHandler(t *testing.T) {
	t.Run("should return 200 and a list of profiles if profile with given id exists", func(t *testing.T) {
		randomId := RandomString()
		randomProfiles := []entities.Profile{RandomProfile(), RandomProfile()}
		followsGetter := func(userId values.UserId) ([]entities.Profile, error) {
			if userId == randomId {
				return randomProfiles, nil
			}
			panic("called with unexpected arguments")
		}

		request := createRequestWithId(randomId)
		response := httptest.NewRecorder()

		handlers.NewGetFollowsHandler(followsGetter).ServeHTTP(response, request)

		randomProfilesResp := handlers.FollowsResponse{Profiles: randomProfiles}
		AssertJSONData(t, response, randomProfilesResp)
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		handlers.NewGetFollowsHandler(nil).ServeHTTP(response, createRequest(nil)) // getter is nil, since it shouldn't be called
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	baseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(userId values.UserId) ([]entities.Profile, error) {
			return nil, err
		}
		handlers.NewGetFollowsHandler(getter).ServeHTTP(rr, createRequestWithId("42"))
	})
}

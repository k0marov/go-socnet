package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0marov/socnet/features/profiles/delivery/http/handlers"
	"github.com/k0marov/socnet/features/profiles/domain/entities"

	"github.com/k0marov/socnet/core/client_errors"
	core_entities "github.com/k0marov/socnet/core/core_entities"
	"github.com/k0marov/socnet/core/core_values"
	helpers "github.com/k0marov/socnet/core/http_test_helpers"
	. "github.com/k0marov/socnet/core/test_helpers"

	"github.com/go-chi/chi/v5"
)

func TestGetMeHandler(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	createRequestWithAuth := func() *http.Request {
		return helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), authUser)
	}
	helpers.BaseTest401(t, handlers.NewGetMeHandler(nil))
	t.Run("should return 200 and a profile if authentication details are provided via context", func(t *testing.T) {
		wantedProfile := RandomContextedProfile()
		getter := func(gotUser, caller core_values.UserId) (entities.ContextedProfile, error) {
			if gotUser == user.Id && caller == user.Id {
				return wantedProfile, nil
			}
			panic(fmt.Sprintf("called with user=%v", gotUser))
		}

		response := httptest.NewRecorder()
		handlers.NewGetMeHandler(getter).ServeHTTP(response, createRequestWithAuth())

		AssertJSONData(t, response, handlers.EntityToResponse(wantedProfile))
	})
	helpers.BaseTestServiceErrorHandling(t, func(wantErr error, response *httptest.ResponseRecorder) {
		getter := func(core_values.UserId, core_values.UserId) (entities.ContextedProfile, error) {
			return entities.ContextedProfile{}, wantErr
		}
		handlers.NewGetMeHandler(getter).ServeHTTP(response, createRequestWithAuth())
	})
}
func createRequestWithId(userId core_values.UserId) *http.Request {
	request := helpers.CreateRequest(nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", userId)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
	return request
}

func TestGetByIdHandler(t *testing.T) {
	helpers.BaseTest401(t, handlers.NewGetByIdHandler(nil))
	t.Run("happy case", func(t *testing.T) {
		targetId := RandomString()
		caller := RandomAuthUser()
		randomProfile := RandomContextedProfile()
		profileGetter := func(target, callerId core_values.UserId) (entities.ContextedProfile, error) {
			if target == targetId && callerId == caller.Id {
				return randomProfile, nil
			}
			panic("called with unexpected arguments")
		}

		request := helpers.AddAuthDataToRequest(createRequestWithId(targetId), caller)
		response := httptest.NewRecorder()

		handlers.NewGetByIdHandler(profileGetter).ServeHTTP(response, request)
		AssertJSONData(t, response, handlers.EntityToResponse(randomProfile))
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		handlers.NewGetByIdHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(target, caller core_values.UserId) (entities.ContextedProfile, error) {
			return entities.ContextedProfile{}, err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithId(RandomString()), RandomAuthUser())
		handlers.NewGetByIdHandler(getter).ServeHTTP(rr, request)
	})
}

func TestFollowsHandler(t *testing.T) {
	t.Run("should return 200 and a list of profiles if profile with given id exists", func(t *testing.T) {
		randomId := RandomString()
		randomProfiles := []core_values.UserId{RandomString(), RandomString(), RandomString()}
		followsGetter := func(userId core_values.UserId) ([]core_values.UserId, error) {
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
		handlers.NewGetFollowsHandler(nil).ServeHTTP(response, helpers.CreateRequest(nil)) // getter is nil, since it shouldn't be called
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(userId core_values.UserId) ([]core_values.UserId, error) {
			return nil, err
		}
		handlers.NewGetFollowsHandler(getter).ServeHTTP(rr, createRequestWithId("42"))
	})
}

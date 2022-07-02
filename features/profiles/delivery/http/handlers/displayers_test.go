package handlers_test

import (
	"context"
	"fmt"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_entities"
	"github.com/k0marov/go-socnet/core/general/core_values"
	helpers "github.com/k0marov/go-socnet/core/helpers/http_test_helpers"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/delivery/http/responses"

	"github.com/k0marov/go-socnet/features/profiles/delivery/http/handlers"
	"github.com/k0marov/go-socnet/features/profiles/domain/entities"

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

		AssertJSONData(t, response, responses.NewProfileResponse(wantedProfile))
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
		AssertJSONData(t, response, responses.NewProfileResponse(randomProfile))
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
	caller := RandomAuthUser()
	helpers.BaseTest401(t, handlers.NewGetFollowsHandler(nil))
	t.Run("should return 200 and a list of profiles if profile with given id exists", func(t *testing.T) {
		randomId := RandomString()
		randomProfiles := []entities.ContextedProfile{RandomContextedProfile(), RandomContextedProfile()}
		followsGetter := func(userId, callerId core_values.UserId) ([]entities.ContextedProfile, error) {
			if userId == randomId && callerId == caller.Id {
				return randomProfiles, nil
			}
			panic("called with unexpected arguments")
		}

		request := helpers.AddAuthDataToRequest(createRequestWithId(randomId), caller)
		response := httptest.NewRecorder()

		handlers.NewGetFollowsHandler(followsGetter).ServeHTTP(response, request)

		AssertJSONData(t, response, responses.NewProfilesResponse(randomProfiles))
	})
	t.Run("error case - id is not provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := helpers.AddAuthDataToRequest(helpers.CreateRequest(nil), RandomAuthUser())
		handlers.NewGetFollowsHandler(nil).ServeHTTP(response, request)
		AssertClientError(t, response, client_errors.IdNotProvided)
	})
	helpers.BaseTestServiceErrorHandling(t, func(err error, rr *httptest.ResponseRecorder) {
		getter := func(userId, callerId core_values.UserId) ([]entities.ContextedProfile, error) {
			return nil, err
		}
		request := helpers.AddAuthDataToRequest(createRequestWithId("42"), RandomAuthUser())
		handlers.NewGetFollowsHandler(getter).ServeHTTP(rr, request)
	})
}

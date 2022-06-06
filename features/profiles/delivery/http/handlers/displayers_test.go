package handlers_test

import (
	core_entities "core/entities"
	. "core/test_helpers"
	"fmt"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/http/handlers"
	"profiles/domain/entities"
	"testing"
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

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, wantedProfile)
	})
	baseTestServiceErrorHandling(t, func(wantErr error, response *httptest.ResponseRecorder) {
		getter := func(core_entities.User) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, wantErr
		}
		handlers.NewGetMeHandler(getter).ServeHTTP(response, createRequestWithAuth())
	})
}

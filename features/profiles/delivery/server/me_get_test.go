package server_test

import (
	core_entities "core/entities"
	"net/http"
	"net/http/httptest"
	"profiles/delivery/server"
	"profiles/domain/entities"
	"testing"

	. "core/test_helpers"
)

func TestHTTPServer_GET_Me(t *testing.T) {
	authUser := RandomAuthUser()
	user := core_entities.UserFromAuth(authUser)
	createRequest := func() *http.Request {
		return httptest.NewRequest(http.MethodGet, "/profiles/me", nil)
	}
	createRequestWithAuth := func() *http.Request {
		r := createRequest()
		return addAuthDataToRequest(r, authUser)
	}

	baseTest401(t, createRequest)
	t.Run("should return 200 and a profile if authentication details are provided via context", func(t *testing.T) {
		wantedProfile := RandomDetailedProfile()
		service := &StubProfileService{
			getDetailed: func(u core_entities.User) (entities.DetailedProfile, error) {
				if u == user {
					return wantedProfile, nil
				}
				panic("getOrCreate called with improper arguments")
			},
		}
		srv := server.NewHTTPServer(service)

		request := createRequestWithAuth()
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		AssertStatusCode(t, response, http.StatusOK)
		AssertJSONData(t, response, wantedProfile)
	})
	baseTestServerErrorHandling(t, createRequestWithAuth)
}

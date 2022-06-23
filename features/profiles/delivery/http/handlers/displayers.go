package handlers

import (
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"net/http"

	"github.com/k0marov/socnet/features/profiles/domain/service"

	"github.com/k0marov/socnet/core/client_errors"
	helpers "github.com/k0marov/socnet/core/http_helpers"

	"github.com/go-chi/chi/v5"
)

type ProfileResponse struct {
	Id         string
	Username   string
	About      string
	AvatarPath string
	Follows    int
	Followers  int
	IsMine     bool
	IsFollowed bool
}

func EntityToResponse(profile entities.ContextedProfile) ProfileResponse {
	return ProfileResponse{
		Id:         profile.Id,
		Username:   profile.Username,
		About:      profile.About,
		AvatarPath: profile.AvatarPath,
		Follows:    profile.Follows,
		Followers:  profile.Followers,
		IsMine:     profile.IsMine,
		IsFollowed: profile.IsLiked,
	}
}

func NewGetMeHandler(getProfile service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		profile, err := getProfile(user.Id, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, EntityToResponse(profile))
	})
}

func NewGetByIdHandler(profileGetter service.ProfileGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := helpers.GetUserOrAddUnauthorized(w, r)
		if !ok {
			return
		}
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		profile, err := profileGetter(id, user.Id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, EntityToResponse(profile))
	})
}

type FollowsResponse struct {
	Profiles []core_values.UserId
}

func NewGetFollowsHandler(followsGetter service.FollowsGetter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			helpers.ThrowClientError(w, client_errors.IdNotProvided)
			return
		}
		follows, err := followsGetter(id)
		if err != nil {
			helpers.HandleServiceError(w, err)
			return
		}
		helpers.WriteJson(w, FollowsResponse{follows})
	})
}

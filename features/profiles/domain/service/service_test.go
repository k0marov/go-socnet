package service_test

import (
	"fmt"
	"testing"

	"github.com/k0marov/go-socnet/core/static_store"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
	"github.com/k0marov/go-socnet/features/profiles/domain/service"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"

	"github.com/k0marov/go-socnet/core/client_errors"
	"github.com/k0marov/go-socnet/core/core_errors"
	"github.com/k0marov/go-socnet/core/core_values"
	"github.com/k0marov/go-socnet/core/ref"
	. "github.com/k0marov/go-socnet/core/test_helpers"
)

func TestFollowToggler(t *testing.T) {
	testTarget := RandomString()
	testFollower := RandomString()
	t.Run("should forward the call to likeable.LikeToggler with owner set to target", func(t *testing.T) {
		wantErr := RandomError()
		likeToggler := func(target string, owner, caller core_values.UserId) error {
			if target == testTarget && owner == testTarget && caller == testFollower {
				return wantErr
			}
			panic("unexpected args")
		}
		gotErr := service.NewFollowToggler(likeToggler)(testTarget, testFollower)
		AssertError(t, gotErr, wantErr)
	})
}

func TestProfileGetter(t *testing.T) {
	target := RandomString()
	caller := RandomString()
	profile := RandomProfile()
	contextedProfile := RandomContextedProfile()
	getProfile := func(id core_values.UserId) (entities.Profile, error) {
		if id == target {
			return profile, nil
		}
		panic("called with unexpected arguments")
	}
	t.Run("error case - store returns NotFoundErr", func(t *testing.T) {
		getProfile := func(core_values.UserId) (entities.Profile, error) {
			return entities.Profile{}, core_errors.ErrNotFound
		}
		sut := service.NewProfileGetter(getProfile, nil)
		_, err := sut(target, caller)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - store returns some other error", func(t *testing.T) {
		getProfile := func(core_values.UserId) (entities.Profile, error) {
			return entities.Profile{}, RandomError()
		}
		sut := service.NewProfileGetter(getProfile, nil)
		_, err := sut(target, caller)
		AssertSomeError(t, err)
	})
	addContext := func(prof entities.Profile, callerId core_values.UserId) (entities.ContextedProfile, error) {
		if prof == profile && callerId == caller {
			return contextedProfile, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - adding context returns some error", func(t *testing.T) {
		addContext := func(profile2 entities.Profile, id core_values.UserId) (entities.ContextedProfile, error) {
			return entities.ContextedProfile{}, RandomError()
		}
		_, err := service.NewProfileGetter(getProfile, addContext)(target, caller)
		AssertSomeError(t, err)
	})
	t.Run("happy case", func(t *testing.T) {
		gotProfile, err := service.NewProfileGetter(getProfile, addContext)(target, caller)
		AssertNoError(t, err)
		Assert(t, gotProfile, contextedProfile, "returned profile")
	})
}

func TestProfileCreator(t *testing.T) {
	user := RandomUser()
	t.Run("happy case", func(t *testing.T) {
		testNewProfile := models.ProfileModel{
			Id:         user.Id,
			Username:   user.Username,
			About:      service.DefaultAbout,
			AvatarPath: service.DefaultAvatarPath,
		}
		wantCreatedProfile := entities.Profile{
			ProfileModel: testNewProfile,
			Follows:      0,
			Followers:    0,
		}
		storeNew := func(profile models.ProfileModel) error {
			if profile == testNewProfile {
				return nil
			}
			panic(fmt.Sprintf("StoreNew called with unexpected profile: %v", profile))
		}
		sut := service.NewProfileCreator(storeNew)

		gotProfile, err := sut(user)
		AssertNoError(t, err)
		Assert(t, gotProfile, wantCreatedProfile, "created profile")
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		storeNew := func(model models.ProfileModel) error {
			return RandomError()
		}
		sut := service.NewProfileCreator(storeNew)

		_, err := sut(user)

		AssertSomeError(t, err)
	})
}

func TestProfileUpdater(t *testing.T) {
	user := RandomUser()
	testUpdateData := values.ProfileUpdateData{
		About: RandomString(),
	}
	silentValidator := func(values.ProfileUpdateData) (client_errors.ClientError, bool) {
		return client_errors.ClientError{}, true
	}
	t.Run("happy case", func(t *testing.T) {
		wantUpdatedProfile := RandomProfile()
		storeUpdater := func(id string, updData values.ProfileUpdateData) (entities.Profile, error) {
			if id == user.Id && updData == testUpdateData {
				return wantUpdatedProfile, nil
			}
			panic(fmt.Sprintf("update called with unexpected arguments: id: %v and updateData: %v", id, updData))
		}
		sut := service.NewProfileUpdater(silentValidator, storeUpdater)

		gotUpdatedProfile, err := sut(user, testUpdateData)
		AssertNoError(t, err)
		Assert(t, gotUpdatedProfile, wantUpdatedProfile, "the returned profile")
	})
	t.Run("error case - validator throws", func(t *testing.T) {
		clientError := RandomClientError()
		validator := func(updData values.ProfileUpdateData) (client_errors.ClientError, bool) {
			if updData == testUpdateData {
				return clientError, false
			}
			panic(fmt.Sprintf("validator called with unexpected args, updData=%v", updData))
		}
		sut := service.NewProfileUpdater(validator, nil) // store is nil, since it shouldn't be accessed
		_, gotErr := sut(user, testUpdateData)
		AssertError(t, gotErr, clientError)
	})
	t.Run("error case - store throws an error", func(t *testing.T) {
		storeUpdater := func(string, values.ProfileUpdateData) (entities.Profile, error) {
			return entities.Profile{}, RandomError()
		}
		sut := service.NewProfileUpdater(silentValidator, storeUpdater)
		_, err := sut(user, testUpdateData)
		AssertSomeError(t, err)
	})
}

func TestAvatarUpdater(t *testing.T) {
	user := RandomUser()
	data := []byte(RandomString())
	dataRef, _ := ref.NewRef(&data)
	testAvatarData := values.AvatarData{
		Data: dataRef,
	}

	silentValidator := func(values.AvatarData) (client_errors.ClientError, bool) {
		return client_errors.ClientError{}, true
	}

	t.Run("happy case", func(t *testing.T) {
		path := RandomString()
		wantURL := static_store.PathToURL(path)
		storeAvatar := func(userId string, avatarData values.AvatarData) (core_values.FileURL, error) {
			if userId == user.Id && avatarData == testAvatarData {
				return path, nil
			}
			panic(fmt.Sprintf("StoreAvatar called with unexpected arguments: userId=%v and avatarData=%v", userId, avatarData))
		}
		sut := service.NewAvatarUpdater(silentValidator, storeAvatar)

		gotURL, err := sut(user, testAvatarData)
		AssertNoError(t, err)
		Assert(t, gotURL, wantURL, "returned profile")
	})
	t.Run("validator throws", func(t *testing.T) {
		clientError := RandomClientError()
		validator := func(avatar values.AvatarData) (client_errors.ClientError, bool) {
			if avatar == testAvatarData {
				return clientError, false
			}
			panic(fmt.Sprintf("validator called with unexpected args, avatar=%v", avatar))
		}
		sut := service.NewAvatarUpdater(validator, nil) // storeAvatar is nil, since it shouldn't be called

		_, err := sut(user, testAvatarData)
		AssertError(t, err, clientError)
	})
	t.Run("store throws an error", func(t *testing.T) {
		storeAvatar := func(string, values.AvatarData) (core_values.FileURL, error) {
			return "", RandomError()
		}
		sut := service.NewAvatarUpdater(silentValidator, storeAvatar)

		_, err := sut(user, testAvatarData)
		AssertSomeError(t, err)
	})
}

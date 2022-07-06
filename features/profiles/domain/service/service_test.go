package service_test

import (
	"fmt"
	"github.com/k0marov/go-socnet/core/general/client_errors"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/general/core_values/ref"
	"github.com/k0marov/go-socnet/core/general/static_store"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/domain/models"

	"github.com/k0marov/go-socnet/features/profiles/domain/entities"
	"github.com/k0marov/go-socnet/features/profiles/domain/service"
	"github.com/k0marov/go-socnet/features/profiles/domain/values"
)

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
			return entities.Profile{}, core_err.ErrNotFound
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

func TestFollowsGetter(t *testing.T) {
	target := RandomId()
	caller := RandomId()
	follows := []core_values.UserId{RandomId()}
	wantFollows := []entities.ContextedProfile{RandomContextedProfile()}

	getFollows := func(id core_values.UserId) ([]core_values.UserId, error) {
		if id == target {
			return follows, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting follows throws", func(t *testing.T) {
		getFollows := func(core_values.UserId) ([]core_values.UserId, error) {
			return nil, RandomError()
		}
		_, err := service.NewFollowsGetter(getFollows, nil)(target, caller)
		AssertSomeError(t, err)

	})
	getProfile := func(targetId, callerId core_values.UserId) (entities.ContextedProfile, error) {
		if targetId == follows[0] && callerId == caller {
			return wantFollows[0], nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting profile throws", func(t *testing.T) {
		getProfile := func(target, caller core_values.UserId) (entities.ContextedProfile, error) {
			return entities.ContextedProfile{}, RandomError()
		}
		_, err := service.NewFollowsGetter(getFollows, getProfile)(target, caller)
		AssertSomeError(t, err)
	})

	t.Run("happy case", func(t *testing.T) {
		sut := service.NewFollowsGetter(getFollows, getProfile)
		gotFollows, err := sut(target, caller)
		AssertNoError(t, err)
		Assert(t, gotFollows, wantFollows, "returned follows")
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
	validator := func(values.ProfileUpdateData) (client_errors.ClientError, bool) {
		return client_errors.ClientError{}, true
	}
	t.Run("error case - validator throws", func(t *testing.T) {
		clientError := RandomClientError()
		validator := func(updData values.ProfileUpdateData) (client_errors.ClientError, bool) {
			if updData == testUpdateData {
				return clientError, false
			}
			panic(fmt.Sprintf("validator called with unexpected args, updData=%v", updData))
		}
		sut := service.NewProfileUpdater(validator, nil, nil) // store is nil, since it shouldn't be accessed
		_, gotErr := sut(user, testUpdateData)
		AssertError(t, gotErr, clientError)
	})
	storeUpdater := func(id string, updData values.ProfileUpdateData) error {
		if id == user.Id && updData == testUpdateData {
			return nil
		}
		panic(fmt.Sprintf("update called with unexpected arguments: id: %v and updateData: %v", id, updData))
	}
	t.Run("error case - store throws an error", func(t *testing.T) {
		storeUpdater := func(string, values.ProfileUpdateData) error {
			return RandomError()
		}
		sut := service.NewProfileUpdater(validator, storeUpdater, nil)
		_, err := sut(user, testUpdateData)
		AssertSomeError(t, err)
	})
	wantUpdatedProfile := RandomContextedProfile()
	profileGetter := func(target, caller core_values.UserId) (entities.ContextedProfile, error) {
		if target == user.Id && caller == user.Id {
			return wantUpdatedProfile, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting updated profile throws", func(t *testing.T) {
		tErr := RandomError()
		profileGetter := func(target, caller core_values.UserId) (entities.ContextedProfile, error) {
			return entities.ContextedProfile{}, tErr
		}
		_, err := service.NewProfileUpdater(validator, storeUpdater, profileGetter)(user, testUpdateData)
		AssertError(t, err, tErr)
	})
	t.Run("happy case", func(t *testing.T) {
		sut := service.NewProfileUpdater(validator, storeUpdater, profileGetter)

		gotUpdatedProfile, err := sut(user, testUpdateData)
		AssertNoError(t, err)
		Assert(t, gotUpdatedProfile, wantUpdatedProfile, "the returned profile")
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

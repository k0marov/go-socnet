package service_test

import (
	"core/client_errors"
	"core/core_errors"
	. "core/test_helpers"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/service"
	"profiles/domain/values"
	"testing"
)

func TestDetailedProfileGetter(t *testing.T) {
	user := RandomUser()
	t.Run("happy case", func(t *testing.T) {
		wantProfile := RandomDetailedProfile()
		storeGetter := func(userId string) (entities.DetailedProfile, error) {
			if userId == user.Id {
				return wantProfile, nil
			}
			panic("GetById called with incorrect arguments")
		}
		sut := service.NewDetailedProfileGetter(storeGetter)

		gotProfile, err := sut(user)

		AssertNoError(t, err)
		Assert(t, gotProfile, wantProfile, "returned profile")
	})
	t.Run("error case - profile does not exist", func(t *testing.T) {
		storeGetter := func(string) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, core_errors.ErrNotFound
		}
		sut := service.NewDetailedProfileGetter(storeGetter)
		_, err := sut(user)
		AssertError(t, err, client_errors.ProfileNotFound)
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		storeGetter := func(userId string) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, RandomError()
		}
		sut := service.NewDetailedProfileGetter(storeGetter)
		_, err := sut(user)
		AssertSomeError(t, err)
	})
}

func TestProfileCreator(t *testing.T) {
	user := RandomUser()
	t.Run("happy case", func(t *testing.T) {
		wantProfile := entities.DetailedProfile{
			Profile: entities.Profile{
				Id:         user.Id,
				Username:   user.Username,
				About:      service.DefaultAbout,
				AvatarPath: service.DefaultAvatarPath,
			},
		}
		storeNew := func(profile entities.DetailedProfile) error {
			if profile == wantProfile {
				return nil
			}
			panic(fmt.Sprintf("StoreNew called with unexpected profile: %v", profile))
		}
		sut := service.NewProfileCreator(storeNew)

		gotProfile, err := sut(user)
		AssertNoError(t, err)
		Assert(t, gotProfile, wantProfile, "created profile")
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		storeNew := func(entities.DetailedProfile) error {
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
		wantUpdatedProfile := RandomDetailedProfile()
		storeUpdater := func(id string, updData values.ProfileUpdateData) (entities.DetailedProfile, error) {
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
		storeUpdater := func(string, values.ProfileUpdateData) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, RandomError()
		}
		sut := service.NewProfileUpdater(silentValidator, storeUpdater)
		_, err := sut(user, testUpdateData)
		AssertSomeError(t, err)
	})
}

func TestAvatarUpdater(t *testing.T) {
	user := RandomUser()
	data := []byte(RandomString())
	testAvatarData := values.AvatarData{
		Data:     &data,
		FileName: RandomString(),
	}

	silentValidator := func(values.AvatarData) (client_errors.ClientError, bool) {
		return client_errors.ClientError{}, true
	}

	t.Run("happy case", func(t *testing.T) {
		wantURL := values.AvatarURL{Url: RandomString()}
		storeAvatar := func(userId string, avatarData values.AvatarData) (values.AvatarURL, error) {
			if userId == user.Id && avatarData == testAvatarData {
				return wantURL, nil
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
		storeAvatar := func(string, values.AvatarData) (values.AvatarURL, error) {
			return values.AvatarURL{}, RandomError()
		}
		sut := service.NewAvatarUpdater(silentValidator, storeAvatar)

		_, err := sut(user, testAvatarData)
		AssertSomeError(t, err)
	})
}

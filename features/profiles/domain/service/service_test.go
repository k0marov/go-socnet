package service_test

import (
	"fmt"
	"testing"

	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/service"
	"github.com/k0marov/socnet/features/profiles/domain/values"

	"github.com/k0marov/socnet/core/client_errors"
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/ref"
	. "github.com/k0marov/socnet/core/test_helpers"
)

func TestFollowsGetter(t *testing.T) {
	userId := RandomString()
	t.Run("happy case", func(t *testing.T) {
		randomFollows := []core_values.UserId{RandomString(), RandomString()}
		storeFollowsGetter := func(gotUserId core_values.UserId) ([]core_values.UserId, error) {
			if gotUserId == userId {
				return randomFollows, nil
			}
			panic("called with unexpected args")
		}
		sut := service.NewFollowsGetter(storeFollowsGetter)

		gotFollows, err := sut(userId)
		AssertNoError(t, err)
		Assert(t, gotFollows, randomFollows, "returned follows")
	})
	t.Run("error case - store returns not found", func(t *testing.T) {
		storeFollowsGetter := func(core_values.UserId) ([]core_values.UserId, error) {
			return nil, core_errors.ErrNotFound
		}
		sut := service.NewFollowsGetter(storeFollowsGetter)
		_, err := sut(userId)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - store returns some other error", func(t *testing.T) {
		storeFollowsGetter := func(core_values.UserId) ([]core_values.UserId, error) {
			return nil, RandomError()
		}
		sut := service.NewFollowsGetter(storeFollowsGetter)
		_, err := sut(userId)
		AssertSomeError(t, err)
	})
}

func TestFollowToggler(t *testing.T) {
	testTarget := RandomString()
	testFollower := RandomString()
	t.Run("should return client error if trying to follow yourself", func(t *testing.T) {
		sut := service.NewFollowToggler(nil, nil, nil)
		err := sut("42", "42")
		AssertError(t, err, client_errors.FollowingYourself)
	})
	t.Run("checking if target is already followed", func(t *testing.T) {
		t.Run("target does not exist", func(t *testing.T) {
			followChecker := func(target, follower core_values.UserId) (bool, error) {
				if target == testTarget && follower == testFollower {
					return false, core_errors.ErrNotFound
				}
				panic("called with unexpected args")
			}
			sut := service.NewFollowToggler(followChecker, nil, nil)
			err := sut(testTarget, testFollower)
			AssertError(t, err, client_errors.NotFound)
		})
		t.Run("some other error is returned", func(t *testing.T) {
			followChecker := func(target, follower core_values.UserId) (bool, error) {
				return false, RandomError()
			}
			sut := service.NewFollowToggler(followChecker, nil, nil)
			err := sut(testTarget, testFollower)
			AssertSomeError(t, err)
		})
	})
	t.Run("target is already followed - unfollow it", func(t *testing.T) {
		followChecker := func(target, follower core_values.UserId) (bool, error) {
			return true, nil
		}
		t.Run("happy case", func(t *testing.T) {
			storeUnfollower := func(target, follower core_values.UserId) error {
				return nil
			}
			sut := service.NewFollowToggler(followChecker, nil, storeUnfollower)
			err := sut(testTarget, testFollower)
			AssertNoError(t, err)
		})
		t.Run("error case - store throws", func(t *testing.T) {
			storeUnfollower := func(target, follower core_values.UserId) error {
				return RandomError()
			}
			sut := service.NewFollowToggler(followChecker, nil, storeUnfollower)
			err := sut(testTarget, testFollower)
			AssertSomeError(t, err)
		})
	})
	t.Run("target is not already followed - follow it", func(t *testing.T) {
		followChecker := func(target, follower core_values.UserId) (bool, error) {
			return false, nil
		}
		t.Run("happy case", func(t *testing.T) {
			storeFollower := func(target, follower core_values.UserId) error {
				return nil
			}
			sut := service.NewFollowToggler(followChecker, storeFollower, nil)
			err := sut(testTarget, testFollower)
			AssertNoError(t, err)
		})
		t.Run("error case - store throws", func(t *testing.T) {
			storeFollower := func(target, follower core_values.UserId) error {
				return RandomError()
			}
			sut := service.NewFollowToggler(followChecker, storeFollower, nil)
			err := sut(testTarget, testFollower)
			AssertSomeError(t, err)
		})
	})
}

func TestProfileGetter(t *testing.T) {
	targetId := RandomString()
	callerId := RandomString()
	t.Run("happy case", func(t *testing.T) {
		randomProfile := RandomProfile()
		isFollowed := RandomBool()
		wantContextedProfile := entities.ContextedProfile{
			Profile:            randomProfile,
			IsFollowedByCaller: isFollowed,
		}
		storeProfileGetter := func(gotUserId core_values.UserId) (entities.Profile, error) {
			if gotUserId == targetId {
				return randomProfile, nil
			}
			panic("called with unexpected arguments")
		}
		storeFollowsChecker := func(target, follower core_values.UserId) (bool, error) {
			if target == targetId && follower == callerId {
				return isFollowed, nil
			}
			panic("unexpected args")
		}
		sut := service.NewProfileGetter(storeProfileGetter, storeFollowsChecker)

		gotProfile, err := sut(targetId, callerId)
		AssertNoError(t, err)
		Assert(t, gotProfile, wantContextedProfile, "returned profile")
	})
	t.Run("error case - store returns NotFoundErr", func(t *testing.T) {
		storeProfileGetter := func(core_values.UserId) (entities.Profile, error) {
			return entities.Profile{}, core_errors.ErrNotFound
		}
		sut := service.NewProfileGetter(storeProfileGetter, nil)
		_, err := sut(targetId, callerId)
		AssertError(t, err, client_errors.NotFound)
	})
	t.Run("error case - store returns some other error", func(t *testing.T) {
		storeProfileGetter := func(core_values.UserId) (entities.Profile, error) {
			return entities.Profile{}, RandomError()
		}
		sut := service.NewProfileGetter(storeProfileGetter, nil)
		_, err := sut(targetId, callerId)
		AssertSomeError(t, err)
	})
	t.Run("error case - store returns error when getting FollowsChecker", func(t *testing.T) {
		storeProfileGetter := func(core_values.UserId) (entities.Profile, error) {
			return entities.Profile{}, nil
		}
		storeFollowsChecker := func(core_values.UserId, core_values.UserId) (bool, error) {
			return false, RandomError()
		}
		_, err := service.NewProfileGetter(storeProfileGetter, storeFollowsChecker)(targetId, callerId)
		AssertSomeError(t, err)
	})
}

func TestProfileCreator(t *testing.T) {
	user := RandomUser()
	t.Run("happy case", func(t *testing.T) {
		testNewProfile := values.NewProfile{
			Id:         user.Id,
			Username:   user.Username,
			About:      service.DefaultAbout,
			AvatarPath: service.DefaultAvatarPath,
		}
		wantCreatedProfile := entities.Profile{
			Id:         user.Id,
			Username:   user.Username,
			About:      service.DefaultAbout,
			AvatarPath: service.DefaultAvatarPath,
			Follows:    0,
			Followers:  0,
		}
		storeNew := func(profile values.NewProfile) error {
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
		storeNew := func(values.NewProfile) error {
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
		wantURL := RandomString()
		storeAvatar := func(userId string, avatarData values.AvatarData) (core_values.ImageUrl, error) {
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
		storeAvatar := func(string, values.AvatarData) (core_values.ImageUrl, error) {
			return "", RandomError()
		}
		sut := service.NewAvatarUpdater(silentValidator, storeAvatar)

		_, err := sut(user, testAvatarData)
		AssertSomeError(t, err)
	})
}

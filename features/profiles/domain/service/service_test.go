package service_test

import (
	"core/client_errors"
	. "core/test_helpers"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/service"
	"testing"
)

func TestService_GetDetailed(t *testing.T) {
	user := RandomUser()
	t.Run("happy case", func(t *testing.T) {
		wantProfile := RandomDetailedProfile()
		store := &MockProfileStore{
			getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
				if userId == user.Id {
					return wantProfile, nil
				}
				panic("GetById called with incorrect arguments")
			},
		}
		sut := service.NewProfileService(store)

		gotProfile, err := sut.GetDetailed(user)

		AssertNoError(t, err)
		Assert(t, gotProfile, wantProfile, "returned profile")
	})
	t.Run("error case - profile does not exist", func(t *testing.T) {
		store := &MockProfileStore{
			getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
				return entities.DetailedProfile{}, service.ErrProfileNotFound
			},
		}
		sut := service.NewProfileService(store)

		_, err := sut.GetDetailed(user)

		AssertError(t, err, client_errors.ProfileNotFound)
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		store := &MockProfileStore{
			getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
				return entities.DetailedProfile{}, RandomError()
			},
		}
		sut := service.NewProfileService(store)

		_, err := sut.GetDetailed(user)
		AssertSomeError(t, err)
	})
}
func TestService_CreateProfileForUser(t *testing.T) {
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
		store := &MockProfileStore{
			storeNew: func(profile entities.DetailedProfile) error {
				if profile == wantProfile {
					return nil
				}
				panic(fmt.Sprintf("StoreNew called with unexpected profile: %v", profile))
			},
		}
		sut := service.NewProfileService(store)

		gotProfile, err := sut.CreateProfileForUser(user)

		AssertNoError(t, err)
		Assert(t, gotProfile, wantProfile, "created profile")
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		store := &MockProfileStore{
			storeNew: func(entities.DetailedProfile) error {
				return RandomError()
			},
		}
		sut := service.NewProfileService(store)

		_, err := sut.CreateProfileForUser(user)

		AssertSomeError(t, err)
	})
}

type MockProfileStore struct {
	getByIdDetailed func(string) (entities.DetailedProfile, error)
	storeNew        func(entities.DetailedProfile) error
}

func (s *MockProfileStore) GetByIdDetailed(userId string) (entities.DetailedProfile, error) {
	if s.getByIdDetailed != nil {
		return s.getByIdDetailed(userId)
	}
	panic("GetById shouldn't have been called")
}

func (s *MockProfileStore) StoreNew(newProfile entities.DetailedProfile) error {
	if s.storeNew != nil {
		return s.storeNew(newProfile)
	}
	panic("CreateNew shouldn't have been called")
}

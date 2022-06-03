package service_test

import (
	. "core/test_helpers"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/service"
	"testing"
)

func TestService_GetOrCreateDetailed(t *testing.T) {
	user := RandomUser()
	t.Run("profile already exists (get it)", func(t *testing.T) {
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

			gotProfile, err := sut.GetOrCreateDetailed(user)

			AssertNoError(t, err)
			Assert(t, gotProfile, wantProfile, "returned profile")
		})
		t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
			store := &MockProfileStore{
				getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
					if userId == user.Id {
						return entities.DetailedProfile{}, RandomError()
					}
					panic(fmt.Sprintf("GetById called with unexpected id: %s", userId))
				},
			}
			sut := service.NewProfileService(store)

			_, err := sut.GetOrCreateDetailed(user)
			AssertSomeError(t, err)
		})
	})
	t.Run("profile does not yet exist (create and return it)", func(t *testing.T) {
		t.Run("happy case", func(t *testing.T) {
			wantProfile := entities.DetailedProfile{
				Profile: entities.Profile{
					Id:       user.Id,
					Username: user.Username,
					About:    service.DefaultAbout,
				},
			}
			store := &MockProfileStore{
				getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
					if userId == user.Id {
						return entities.DetailedProfile{}, service.ErrProfileNotFound
					}
					panic(fmt.Sprintf("GetById called with unexpected id: %s", userId))
				},
				storeNew: func(profile entities.DetailedProfile) error {
					if profile == wantProfile {
						return nil
					}
					panic(fmt.Sprintf("StoreNew called with unexpected profile: %v", profile))
				},
			}
			sut := service.NewProfileService(store)

			gotProfile, err := sut.GetOrCreateDetailed(user)

			AssertNoError(t, err)
			Assert(t, gotProfile, wantProfile, "returned profile")
		})
		t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
			t.Run("store throws when getting a profile", func(t *testing.T) {
				store := &MockProfileStore{
					getByIdDetailed: func(s string) (entities.DetailedProfile, error) {
						return entities.DetailedProfile{}, RandomError()
					},
				}
				sut := service.NewProfileService(store)

				_, err := sut.GetOrCreateDetailed(user)

				AssertSomeError(t, err)
			})
			t.Run("store throws when creating a profile", func(t *testing.T) {
				store := &MockProfileStore{
					getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
						if userId == user.Id {
							return entities.DetailedProfile{}, service.ErrProfileNotFound
						}
						panic(fmt.Sprintf("GetById called with unexpected id: %s", userId))
					},
					storeNew: func(entities.DetailedProfile) error {
						return RandomError()
					},
				}
				sut := service.NewProfileService(store)

				_, err := sut.GetOrCreateDetailed(user)

				AssertSomeError(t, err)
			})
		})
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

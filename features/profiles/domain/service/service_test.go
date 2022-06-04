package service_test

import (
	"core/client_errors"
	"core/image_decoder"
	. "core/test_helpers"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"profiles/domain/entities"
	"profiles/domain/service"
	"profiles/domain/values"
	"strings"
	"testing"
)

var panickingImageDecoder = PanickingImageDecoder{}

func TestService_Update(t *testing.T) {
	user := RandomUser()
	t.Run("field validation", func(t *testing.T) {
		updatedProfile := RandomDetailedProfile()
		tooLongAbout := strings.Repeat("Loooooong about!", 100)
		successStore := &MockProfileStore{
			update: func(id string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
				return updatedProfile, nil
			},
		}
		panickingStore := &MockProfileStore{}
		cases := []struct {
			updateData    values.ProfileUpdateData
			expectedError error
		}{
			{values.ProfileUpdateData{About: RandomString()}, nil},
			{values.ProfileUpdateData{About: tooLongAbout}, client_errors.AboutTooLong},
		}
		for _, c := range cases {
			var store *MockProfileStore
			if c.expectedError != nil {
				store = panickingStore
			} else {
				store = successStore
			}
			sut := service.NewProfileService(store, panickingImageDecoder)

			_, err := sut.Update(user, c.updateData)
			Assert(t, err, c.expectedError, "returned error")
		}
	})
	t.Run("should update the profile with proper arguments using store if all checks have passed", func(t *testing.T) {
		goodUpdateData := values.ProfileUpdateData{
			About: RandomString(),
		}
		t.Run("happy case", func(t *testing.T) {
			wantUpdatedProfile := RandomDetailedProfile()
			store := &MockProfileStore{
				update: func(id string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
					if id == user.Id && updateData == goodUpdateData {
						return wantUpdatedProfile, nil
					}
					panic(fmt.Sprintf("update called with unexpected arguments: id: %v and updateData: %v", id, updateData))
				},
			}
			sut := service.NewProfileService(store, panickingImageDecoder)

			gotUpdatedProfile, err := sut.Update(user, goodUpdateData)
			AssertNoError(t, err)
			Assert(t, gotUpdatedProfile, wantUpdatedProfile, "the returned profile")
		})
		t.Run("error case - store throws an error", func(t *testing.T) {
			t.Run("it is a not found error", func(t *testing.T) {
				store := &MockProfileStore{
					update: func(s string, pud values.ProfileUpdateData) (entities.DetailedProfile, error) {
						return entities.DetailedProfile{}, service.ErrProfileNotFound
					},
				}
				sut := service.NewProfileService(store, panickingImageDecoder)

				_, err := sut.Update(user, goodUpdateData)
				AssertError(t, err, client_errors.ProfileNotFound)
			})
			t.Run("it is some other error", func(t *testing.T) {
				store := &MockProfileStore{
					update: func(string, values.ProfileUpdateData) (entities.DetailedProfile, error) {
						return entities.DetailedProfile{}, RandomError()
					},
				}
				sut := service.NewProfileService(store, panickingImageDecoder)

				_, err := sut.Update(user, goodUpdateData)
				AssertSomeError(t, err)
			})
		})
	})
}

func TestService_UpdateAvatar(t *testing.T) {
	user := RandomUser()
	realImageDecoder := image_decoder.ImageDecoderImpl{}
	readFixtureFile := func(fileName string) *[]byte {
		file, _ := os.Open(filepath.Join("testdata", fileName))
		fileBytes, _ := io.ReadAll(file)
		return &fileBytes
	}
	goodAvatarFileName := "test_avatar.jpg"

	t.Run("avatar file validation", func(t *testing.T) {
		panickingStore := &MockProfileStore{}
		silentStore := &MockProfileStore{
			storeAvatar: func(s string, ad values.AvatarData) (entities.DetailedProfile, error) {
				return RandomDetailedProfile(), nil
			},
		}
		cases := []struct {
			fixtureFilename string
			expectedError   error
		}{
			{goodAvatarFileName, nil},
			{"test_text.txt", client_errors.NonImageAvatar},
			{"test_js_injection.js", client_errors.NonImageAvatar},
			{"test_non_square_avatar.png", client_errors.NonSquareAvatar},
		}
		for _, c := range cases {
			t.Run(c.fixtureFilename, func(t *testing.T) {
				fileBytes := readFixtureFile(c.fixtureFilename)
				avatarData := values.AvatarData{
					Data:     fileBytes,
					FileName: c.fixtureFilename,
				}
				var store *MockProfileStore
				if c.expectedError != nil {
					store = panickingStore
				} else {
					store = silentStore
				}
				sut := service.NewProfileService(store, realImageDecoder)

				_, err := sut.UpdateAvatar(user, avatarData)
				Assert(t, err, c.expectedError, "returned error")
			})
		}
	})
	t.Run("if validation has passed, should call store with proper arguments", func(t *testing.T) {
		getAvatarData := func() values.AvatarData {
			fileBytes := readFixtureFile(goodAvatarFileName)
			return values.AvatarData{
				Data:     fileBytes,
				FileName: RandomString(),
			}
		}
		t.Run("happy case", func(t *testing.T) {
			goodAvatarData := getAvatarData()
			wantUpdatedProfile := RandomDetailedProfile()
			store := &MockProfileStore{
				storeAvatar: func(userId string, avatarData values.AvatarData) (entities.DetailedProfile, error) {
					if userId == user.Id && avatarData == goodAvatarData {
						return wantUpdatedProfile, nil
					}
					panic(fmt.Sprintf("StoreAvatar called with unexpected arguments: userId=%v and avatarData=%v", userId, avatarData))
				},
			}
			sut := service.NewProfileService(store, realImageDecoder)

			gotUpdatedProfile, err := sut.UpdateAvatar(user, goodAvatarData)
			AssertNoError(t, err)
			Assert(t, gotUpdatedProfile, wantUpdatedProfile, "returned profile")
		})
		t.Run("store throws an error", func(t *testing.T) {
			goodAvatarData := getAvatarData()
			store := &MockProfileStore{
				storeAvatar: func(string, values.AvatarData) (entities.DetailedProfile, error) {
					return entities.DetailedProfile{}, RandomError()
				},
			}
			sut := service.NewProfileService(store, realImageDecoder)

			_, err := sut.UpdateAvatar(user, goodAvatarData)
			AssertSomeError(t, err)
		})
	})
}

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
		sut := service.NewProfileService(store, panickingImageDecoder)

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
		sut := service.NewProfileService(store, panickingImageDecoder)

		_, err := sut.GetDetailed(user)

		AssertError(t, err, client_errors.ProfileNotFound)
	})
	t.Run("error case - store throws, it is NOT a client error", func(t *testing.T) {
		store := &MockProfileStore{
			getByIdDetailed: func(userId string) (entities.DetailedProfile, error) {
				return entities.DetailedProfile{}, RandomError()
			},
		}
		sut := service.NewProfileService(store, panickingImageDecoder)

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
		sut := service.NewProfileService(store, panickingImageDecoder)

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
		sut := service.NewProfileService(store, panickingImageDecoder)

		_, err := sut.CreateProfileForUser(user)

		AssertSomeError(t, err)
	})
}

type MockProfileStore struct {
	getByIdDetailed func(string) (entities.DetailedProfile, error)
	storeNew        func(entities.DetailedProfile) error
	update          func(string, values.ProfileUpdateData) (entities.DetailedProfile, error)
	storeAvatar     func(string, values.AvatarData) (entities.DetailedProfile, error)
}

func (m *MockProfileStore) GetByIdDetailed(userId string) (entities.DetailedProfile, error) {
	if m.getByIdDetailed != nil {
		return m.getByIdDetailed(userId)
	}
	panic("GetById shouldn't have been called")
}

func (m *MockProfileStore) StoreNew(newProfile entities.DetailedProfile) error {
	if m.storeNew != nil {
		return m.storeNew(newProfile)
	}
	panic("CreateNew shouldn't have been called")
}

func (m *MockProfileStore) Update(userId string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
	if m.update != nil {
		return m.update(userId, updateData)
	}
	panic("Update shouldn't have been called")
}

func (m *MockProfileStore) StoreAvatar(userId string, avatar values.AvatarData) (entities.DetailedProfile, error) {
	if m.storeAvatar != nil {
		return m.storeAvatar(userId, avatar)
	}
	panic("StoreAvatar shouldn't have been called")
}

type PanickingImageDecoder struct{}

func (m PanickingImageDecoder) Decode(*[]byte) (image_decoder.Image, error) {
	panic("Decode shouldn't have been called")
}

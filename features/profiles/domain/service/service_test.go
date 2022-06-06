package service_test

import (
	"core/client_errors"
	"core/core_errors"
	"core/image_decoder"
	. "core/test_helpers"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/service"
	"profiles/domain/values"
	"strings"
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
	goodUpdateData := values.ProfileUpdateData{
		About: RandomString(),
	}
	t.Run("field validation", func(t *testing.T) {
		goodStoreUpdater := func(string, values.ProfileUpdateData) (entities.DetailedProfile, error) {
			return RandomDetailedProfile(), nil
		}
		cases := []struct {
			updateData    values.ProfileUpdateData
			expectedError error
		}{
			{goodUpdateData, nil},
			{values.ProfileUpdateData{About: strings.Repeat("long", 100)}, client_errors.AboutTooLong},
		}
		for _, c := range cases {
			var storeUpdater service.StoreProfileUpdater
			if c.expectedError == nil {
				storeUpdater = goodStoreUpdater
			}
			sut := service.NewProfileUpdater(storeUpdater) // if error != nil, updater shouldn't be called, so it's nil
			_, err := sut(user, c.updateData)
			Assert(t, err, c.expectedError, "returned error")
		}
	})
	t.Run("should update the profile with proper arguments using store if all checks have passed", func(t *testing.T) {
		t.Run("happy case", func(t *testing.T) {
			wantUpdatedProfile := RandomDetailedProfile()
			storeUpdater := func(id string, updateData values.ProfileUpdateData) (entities.DetailedProfile, error) {
				if id == user.Id && updateData == goodUpdateData {
					return wantUpdatedProfile, nil
				}
				panic(fmt.Sprintf("update called with unexpected arguments: id: %v and updateData: %v", id, updateData))
			}
			sut := service.NewProfileUpdater(storeUpdater)

			gotUpdatedProfile, err := sut(user, goodUpdateData)
			AssertNoError(t, err)
			Assert(t, gotUpdatedProfile, wantUpdatedProfile, "the returned profile")
		})
		t.Run("error case - store throws an error", func(t *testing.T) {
			t.Run("it is a not found error", func(t *testing.T) {
				storeUpdater := func(string, values.ProfileUpdateData) (entities.DetailedProfile, error) {
					return entities.DetailedProfile{}, core_errors.ErrNotFound
				}
				sut := service.NewProfileUpdater(storeUpdater)
				_, err := sut(user, goodUpdateData)
				AssertError(t, err, client_errors.ProfileNotFound)
			})
			t.Run("it is some other error", func(t *testing.T) {
				storeUpdater := func(string, values.ProfileUpdateData) (entities.DetailedProfile, error) {
					return entities.DetailedProfile{}, RandomError()
				}
				sut := service.NewProfileUpdater(storeUpdater)
				_, err := sut(user, goodUpdateData)
				AssertSomeError(t, err)
			})
		})
	})
}

func TestAvatarUpdater(t *testing.T) {
	user := RandomUser()
	goodAvatar := []byte(RandomString())
	jsAvatar := []byte(RandomString())
	nonSquareAvatar := []byte(RandomString())

	stubImageDecoder := func(avatar *[]byte) (image_decoder.Image, error) {
		if avatar == &goodAvatar {
			return image_decoder.Image{Width: 10, Height: 10}, nil
		} else if avatar == &nonSquareAvatar {
			return image_decoder.Image{Width: 42, Height: 24}, nil
		} else if avatar == &jsAvatar {
			return image_decoder.Image{}, RandomError()
		}
		panic(fmt.Sprintf("image decoder called with unexpected avatar=%v", avatar))
	}

	t.Run("avatar file validation", func(t *testing.T) {
		stubStoreAvatar := func(s string, ad values.AvatarData) (values.AvatarURL, error) {
			return values.AvatarURL{}, nil
		}
		// stubImageDecoder := func()
		cases := []struct {
			name          string
			avatarBytes   *[]byte
			expectedError error
		}{
			{"Happy case", &goodAvatar, nil},
			{"Avatar not really an image", &jsAvatar, client_errors.NonImageAvatar},
			{"Avatar not square", &nonSquareAvatar, client_errors.NonSquareAvatar},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				avatarData := values.AvatarData{Data: c.avatarBytes, FileName: RandomString()}
				var storeAvatar service.StoreAvatarUpdater
				if c.expectedError == nil {
					storeAvatar = stubStoreAvatar
				} // else storeAvatar shouldn't be called, so let it be nil

				sut := service.NewAvatarUpdater(storeAvatar, stubImageDecoder)

				_, err := sut(user, avatarData)
				Assert(t, err, c.expectedError, "returned error")
			})
		}
	})
	t.Run("if validation has passed, should call store with proper arguments", func(t *testing.T) {
		data := []byte(RandomString())
		goodAvatarData := values.AvatarData{
			Data:     &data,
			FileName: RandomString(),
		}
		silentImageDecoder := func(*[]byte) (image_decoder.Image, error) {
			return image_decoder.Image{Width: 10, Height: 10}, nil
		}
		t.Run("happy case", func(t *testing.T) {
			wantURL := values.AvatarURL{Url: RandomString()}
			storeAvatar := func(userId string, avatarData values.AvatarData) (values.AvatarURL, error) {
				if userId == user.Id && avatarData == goodAvatarData {
					return wantURL, nil
				}
				panic(fmt.Sprintf("StoreAvatar called with unexpected arguments: userId=%v and avatarData=%v", userId, avatarData))
			}
			sut := service.NewAvatarUpdater(storeAvatar, silentImageDecoder)

			gotURL, err := sut(user, goodAvatarData)
			AssertNoError(t, err)
			Assert(t, gotURL, wantURL, "returned profile")
		})
		t.Run("store throws an error", func(t *testing.T) {
			storeAvatar := func(string, values.AvatarData) (values.AvatarURL, error) {
				return values.AvatarURL{}, RandomError()
			}
			sut := service.NewAvatarUpdater(storeAvatar, silentImageDecoder)

			_, err := sut(user, goodAvatarData)
			AssertSomeError(t, err)
		})
	})
}

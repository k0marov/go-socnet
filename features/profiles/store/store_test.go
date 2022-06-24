package store_test

import (
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/static_store"
	"github.com/k0marov/socnet/features/profiles/domain/models"
	"reflect"
	"testing"

	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"
	"github.com/k0marov/socnet/features/profiles/store"

	"github.com/k0marov/socnet/core/ref"
	. "github.com/k0marov/socnet/core/test_helpers"
)

func TestStoreProfileUpdater(t *testing.T) {
	testDetailedProfile := RandomProfile()
	testUpdData := values.ProfileUpdateData{About: RandomString()}
	testDBUpdData := store.DBUpdateData{About: testUpdData.About}
	wantUpdatedProfile := entities.Profile{
		ProfileModel: models.ProfileModel{
			Id:         testDetailedProfile.Id,
			Username:   testDetailedProfile.Username,
			AvatarPath: testDetailedProfile.AvatarPath,
			About:      testUpdData.About,
		},
	}
	t.Run("happy case", func(t *testing.T) {
		updaterCalled := false
		dbUpdater := func(id string, updData store.DBUpdateData) error {
			if id == testDetailedProfile.Id && updData == testDBUpdData {
				updaterCalled = true
				return nil
			}
			panic(fmt.Sprintf("called with unexpected arguments, id=%v, updData=%v", id, updData))
		}
		profileGetter := func(id string) (entities.Profile, error) {
			if id == testDetailedProfile.Id {
				return wantUpdatedProfile, nil
			}
			panic(fmt.Sprintf("called with unexpected arguments, id=%v", id))
		}
		sut := store.NewStoreProfileUpdater(dbUpdater, profileGetter)

		gotProfile, err := sut(testDetailedProfile.Id, testUpdData)
		AssertNoError(t, err)
		Assert(t, updaterCalled, true, "db updater called")
		Assert(t, gotProfile, wantUpdatedProfile, "returned updated profile")
	})
	t.Run("error case - updater returns an error", func(t *testing.T) {
		dbUpdater := func(string, store.DBUpdateData) error {
			return RandomError()
		}
		sut := store.NewStoreProfileUpdater(dbUpdater, nil) // getter is nil, since it shouldn't be called
		_, err := sut(testDetailedProfile.Id, testUpdData)
		AssertSomeError(t, err)
	})
	t.Run("error case - getter returns an error", func(t *testing.T) {
		tErr := RandomError()
		dbUpdater := func(string, store.DBUpdateData) error {
			return nil
		}
		profileGetter := func(string) (entities.Profile, error) {
			return entities.Profile{}, tErr
		}
		sut := store.NewStoreProfileUpdater(dbUpdater, profileGetter)
		_, err := sut(testDetailedProfile.Id, testUpdData)
		AssertError(t, err, tErr)
	})
}

func TestStoreAvatarUpdater(t *testing.T) {
	t.Run("should store avatar using file storage", func(t *testing.T) {
		randomFile := []byte(RandomString())
		randomFileRef, _ := ref.NewRef(&randomFile)
		avatar := values.AvatarData{
			Data: randomFileRef,
		}
		userId := RandomString()
		t.Run("happy case", func(t *testing.T) {
			wantPath := RandomString()
			storeFile := func(file ref.Ref[[]byte], belongsToUser string) (string, error) {
				if reflect.DeepEqual(file, randomFileRef) && belongsToUser == userId {
					return wantPath, nil
				}
				panic(fmt.Sprintf("StoreFile called with unexpected arguments, file=%v, belongsToUser=%v", file, belongsToUser))
			}

			t.Run("should store avatarPath in DB", func(t *testing.T) {
				t.Run("happy case", func(t *testing.T) {
					updateProfile := func(string, store.DBUpdateData) error {
						return nil
					}
					sut := store.NewStoreAvatarUpdater(storeFile, updateProfile)
					gotAvatarUrl, err := sut(userId, avatar)
					AssertNoError(t, err)
					Assert(t, gotAvatarUrl, wantPath, "returned avatar url")
				})
				t.Run("error case - db returns an error", func(t *testing.T) {
					wantUpdData := store.DBUpdateData{AvatarPath: wantPath}
					updateProfile := func(gotUserId string, data store.DBUpdateData) error {
						if gotUserId == userId && data == wantUpdData {
							return RandomError()
						}
						panic(fmt.Sprintf("called with unexpected arguments, gotUserId=%v, avatar=%v", gotUserId, avatar))
					}
					sut := store.NewStoreAvatarUpdater(storeFile, updateProfile)

					_, err := sut(userId, avatar)
					AssertSomeError(t, err)
				})
			})

		})
		t.Run("error case - fileStore returns an error", func(t *testing.T) {
			storeFile := func(ref.Ref[[]byte], string) (string, error) {
				return "", RandomError()
			}
			sut := store.NewStoreAvatarUpdater(storeFile, nil) // nil, because db shouldn't be called

			_, err := sut(userId, avatar)
			AssertSomeError(t, err)
		})
	})
}

func TestStoreProfileGetter(t *testing.T) {
	profileId := RandomId()
	model := RandomProfileModel()
	follows := RandomInt()
	followers := RandomInt()

	dbGetter := func(id core_values.UserId) (models.ProfileModel, error) {
		if id == profileId {
			return model, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting model from db throws", func(t *testing.T) {
		dbGetter := func(core_values.UserId) (models.ProfileModel, error) {
			return models.ProfileModel{}, RandomError()
		}
		_, err := store.NewStoreProfileGetter(dbGetter, nil, nil)(profileId)
		AssertSomeError(t, err)
	})
	followersGetter := func(targetId core_values.UserId) (int, error) {
		if targetId == profileId {
			return followers, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting followers throws", func(t *testing.T) {
		wantErr := RandomError()
		followersGetter := func(core_values.UserId) (int, error) {
			return 0, wantErr
		}
		_, err := store.NewStoreProfileGetter(dbGetter, followersGetter, nil)(profileId)
		AssertError(t, err, wantErr)
	})
	followsGetter := func(targetId core_values.UserId) (int, error) {
		if targetId == profileId {
			return follows, nil
		}
		panic("unexpected args")
	}
	t.Run("error case - getting follows throws", func(t *testing.T) {
		wantErr := RandomError()
		followsGetter := func(id core_values.UserId) (int, error) {
			return 0, wantErr
		}
		_, err := store.NewStoreProfileGetter(dbGetter, followersGetter, followsGetter)(profileId)
		AssertError(t, err, wantErr)
	})

	sut := store.NewStoreProfileGetter(dbGetter, followersGetter, followsGetter)
	gotProfile, err := sut(profileId)
	AssertNoError(t, err)
	wantProfile := entities.Profile{
		ProfileModel: model,
		AvatarURL:    static_store.PathToURL(model.AvatarPath),
		Follows:      follows,
		Followers:    followers,
	}
	Assert(t, gotProfile, wantProfile, "returned profile entity")
}

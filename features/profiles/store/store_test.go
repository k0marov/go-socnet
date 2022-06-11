package store_test

import (
	"core/core_errors"
	"core/ref"
	. "core/test_helpers"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/values"
	"profiles/store"
	"reflect"
	"testing"
)

func TestStoreProfileCreator(t *testing.T) {
	testProfile := RandomProfile()
	t.Run("should forward the call to db", func(t *testing.T) {
		wantErr := RandomError()
		dbCreator := func(gotProfile entities.Profile) error {
			if gotProfile == testProfile {
				return wantErr
			}
			panic(fmt.Sprintf("called with unexpected args, gotProfile=%v", gotProfile))
		}
		sut := store.NewStoreProfileCreator(dbCreator)

		err := sut(testProfile)
		AssertError(t, err, wantErr)

	})
}

func TestStoreDetailedProfileGetter(t *testing.T) { // TODO adding follows in this layer
	t.Run("happy case", func(t *testing.T) {
		wantProfile := RandomProfile()
		wantDetailedProfile := entities.DetailedProfile{Profile: wantProfile}
		dbProfileGetter := func(gotId string) (entities.Profile, error) {
			if gotId == wantProfile.Id {
				return wantProfile, nil
			}
			panic(fmt.Sprintf("called with unexpected args, gotId=%v", wantProfile.Id))
		}
		sut := store.NewStoreDetailedProfileGetter(dbProfileGetter)

		gotDetailedProfile, err := sut(wantProfile.Id)
		AssertNoError(t, err)
		Assert(t, gotDetailedProfile, wantDetailedProfile, "returned detailed profile")
	})
	t.Run("error case - not found error", func(t *testing.T) {
		dbProfileGetter := func(string) (entities.Profile, error) {
			return entities.Profile{}, core_errors.ErrNotFound
		}
		sut := store.NewStoreDetailedProfileGetter(dbProfileGetter)
		_, err := sut(RandomString())
		AssertError(t, err, core_errors.ErrNotFound)
	})
	t.Run("error case - getting profile from db throws oether error", func(t *testing.T) {
		dbProfileGetter := func(string) (entities.Profile, error) {
			return entities.Profile{}, RandomError()
		}
		sut := store.NewStoreDetailedProfileGetter(dbProfileGetter)
		_, err := sut(RandomString())
		AssertSomeError(t, err)
	})
}

func TestStoreProfileUpdater(t *testing.T) {
	testDetailedProfile := RandomDetailedProfile()
	testUpdData := values.ProfileUpdateData{About: RandomString()}
	testDBUpdData := store.DBUpdateData{About: testUpdData.About}
	wantUpdatedProfile := entities.DetailedProfile{
		Profile: entities.Profile{
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
		profileGetter := func(id string) (entities.DetailedProfile, error) {
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
		profileGetter := func(string) (entities.DetailedProfile, error) {
			return entities.DetailedProfile{}, tErr
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
					wantAvatarURL := values.AvatarPath{Path: wantPath}
					updateProfile := func(string, store.DBUpdateData) error {
						return nil
					}
					sut := store.NewStoreAvatarUpdater(storeFile, updateProfile)
					gotAvatarUrl, err := sut(userId, avatar)
					AssertNoError(t, err)
					Assert(t, gotAvatarUrl, wantAvatarURL, "returned avatar url")
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

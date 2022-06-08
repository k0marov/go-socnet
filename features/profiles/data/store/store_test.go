package store_test

import (
	"core/ref"
	. "core/test_helpers"
	"fmt"
	"profiles/data/store"
	"profiles/domain/entities"
	"profiles/domain/values"
	"reflect"
	"testing"
)

func TestStoreProfileCreator(t *testing.T) {
	testProfile := RandomDetailedProfile()
	t.Run("should forward the call to db", func(t *testing.T) {
		wantErr := RandomError()
		dbCreator := func(gotProfile entities.DetailedProfile) error {
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
	t.Run("error case - getting profile from db throws", func(t *testing.T) {
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
		dbUpdater := func(id string, updData values.ProfileUpdateData) error {
			if id == testDetailedProfile.Id && updData == testUpdData {
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
		dbUpdater := func(string, values.ProfileUpdateData) error {
			return RandomError()
		}
		sut := store.NewStoreProfileUpdater(dbUpdater, nil) // getter is nil, since it shouldn't be called
		_, err := sut(testDetailedProfile.Id, testUpdData)
		AssertSomeError(t, err)
	})
	t.Run("error case - getter returns an error", func(t *testing.T) {
		tErr := RandomError()
		dbUpdater := func(string, values.ProfileUpdateData) error {
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
			storeFile := func(file []byte, belongsToUser string) (string, error) {
				if reflect.DeepEqual(file, randomFile) && belongsToUser == userId {
					return wantPath, nil
				}
				panic(fmt.Sprintf("StoreFile called with unexpected arguments, file=%v, belongsToUser=%v", file, belongsToUser))
			}

			t.Run("should store avatarPath in DB", func(t *testing.T) {
				t.Run("happy case", func(t *testing.T) {
					wantAvatarURL := values.AvatarURL{Url: wantPath}
					storeDBAvatar := func(string, values.AvatarURL) error {
						return nil
					}
					sut := store.NewStoreAvatarUpdater(storeFile, storeDBAvatar)
					gotAvatarUrl, err := sut(userId, avatar)
					AssertNoError(t, err)
					Assert(t, gotAvatarUrl, wantAvatarURL, "returned avatar url")
				})
				t.Run("error case - db returns an error", func(t *testing.T) {
					storeDBAvatar := func(gotUserId string, avatar values.AvatarURL) error {
						if gotUserId == userId && avatar.Url == wantPath {
							return RandomError()
						}
						panic(fmt.Sprintf("called with unexpected arguments, gotUserId=%v, avatar=%v", gotUserId, avatar))
					}
					sut := store.NewStoreAvatarUpdater(storeFile, storeDBAvatar)

					_, err := sut(userId, avatar)
					AssertSomeError(t, err)
				})
			})

		})
		t.Run("error case - fileStore returns an error", func(t *testing.T) {
			storeFile := func([]byte, string) (string, error) {
				return "", RandomError()
			}
			sut := store.NewStoreAvatarUpdater(storeFile, nil) // nil, because db shouldn't be called

			_, err := sut(userId, avatar)
			AssertSomeError(t, err)
		})
	})
}

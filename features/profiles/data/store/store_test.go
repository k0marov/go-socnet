package store_test

import (
	. "core/test_helpers"
	"fmt"
	"profiles/data/store"
	"profiles/domain/values"
	"testing"
)

func TestStoreAvatarUpdater(t *testing.T) {
	t.Run("should store avatar using file storage", func(t *testing.T) {
		randomFile := []byte(RandomString())
		avatar := values.AvatarData{
			Data:     &randomFile,
			FileName: RandomString(),
		}
		userId := RandomString()
		t.Run("happy case", func(t *testing.T) {
			wantPath := RandomString()
			storeFile := func(file *[]byte, dir, fileName string) (string, error) {
				if file == &randomFile && dir == store.AvatarsDir && fileName == userId {
					return wantPath, nil
				}
				panic(fmt.Sprintf("StoreFile called with unexpected arguments, file=%v, dir=%v, fileName=%v", file, dir, fileName))
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
			storeFile := func(b *[]byte, s1, s2 string) (string, error) {
				return "", RandomError()
			}
			sut := store.NewStoreAvatarUpdater(storeFile, nil) // nil, because db shouldn't be called

			_, err := sut(userId, avatar)
			AssertSomeError(t, err)
		})
	})
}

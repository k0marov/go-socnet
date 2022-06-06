package store_test

import (
	. "core/test_helpers"
	"fmt"
	"profiles/data/store"
	"profiles/domain/values"
	"testing"
)

type MockFileStorage struct {
	storeFile func(*[]byte, string, string) (string, error)
}

func (m *MockFileStorage) StoreFile(file *[]byte, dir string, fileName string) (string, error) {
	if m.storeFile != nil {
		return m.storeFile(file, dir, fileName)
	}
	panic("StoreFile should not be called here")
}

func TestStore_StoreAvatar(t *testing.T) {
	t.Run("should store avatar using file storage", func(t *testing.T) {
		randomFile := []byte(RandomString())
		avatar := values.AvatarData{
			Data:     &randomFile,
			FileName: RandomString(),
		}
		userId := RandomString()
		t.Run("happy case", func(t *testing.T) {
			wantPath := RandomString()
			storage := &MockFileStorage{
				storeFile: func(file *[]byte, dir, fileName string) (string, error) {
					if file == &randomFile && dir == store.AvatarsDir && fileName == userId {
						return wantPath, nil
					}
					panic(fmt.Sprintf("StoreFile called with unexpected arguments, file=%v, dir=%v, fileName=%v", file, dir, fileName))
				},
			}

			sut := store.NewProfileStoreImpl(storage, nil)

			gotPath, err := sut.StoreAvatar(userId, avatar)
			AssertNoError(t, err)
			Assert(t, gotPath, values.AvatarURL{Url: wantPath}, "the returned path")
		})
		t.Run("error case - store returns an error", func(t *testing.T) {
			storage := &MockFileStorage{
				storeFile: func(b *[]byte, s1, s2 string) (string, error) {
					return "", RandomError()
				},
			}
			sut := store.NewProfileStoreImpl(storage, nil)

			_, err := sut.StoreAvatar(userId, avatar)
			AssertSomeError(t, err)
		})
	})
}

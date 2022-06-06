package store

import (
	"fmt"
	"profiles/domain/values"
)

type FileStorage interface {
	StoreFile(contents *[]byte, directory, fileName string) (path string, err error)
}
type SqlDB interface {
}

type ProfileStoreImpl struct {
	fileStorage FileStorage
	sqlDB       SqlDB
}

func NewProfileStoreImpl(fileStorage FileStorage, sqlDB SqlDB) *ProfileStoreImpl {
	return &ProfileStoreImpl{
		fileStorage: fileStorage,
		sqlDB:       sqlDB,
	}
}

const AvatarsDir = "avatars"

func (p *ProfileStoreImpl) StoreAvatar(userId string, avatar values.AvatarData) (values.AvatarURL, error) {
	path, err := p.fileStorage.StoreFile(avatar.Data, AvatarsDir, userId)
	if err != nil {
		return values.AvatarURL{}, fmt.Errorf("error while storing avatar file: %w", err)
	}
	return values.AvatarURL{Url: path}, nil
}

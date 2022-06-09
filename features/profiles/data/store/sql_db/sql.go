package sql

import (
	"database/sql"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type SqlDB struct {
	sql sql.DB
}

func NewSqlDB(sql sql.DB) *SqlDB {
	return &SqlDB{sql: sql}
}

func (db *SqlDB) UpdateAvatar(userId string, avatarPath values.AvatarURL) error {
	panic("unimplemented")
}

func (db *SqlDB) UpdateProfile(id string, updData values.ProfileUpdateData) error {
	panic("unimplemented")
}

func (db *SqlDB) CreateProfile(entities.DetailedProfile) error {
	panic("unimplemented")
}

func (db *SqlDB) GetProfile(id string) (entities.Profile, error) {
	panic("unimplemented")
}

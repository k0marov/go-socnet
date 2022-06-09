package sql_db

import (
	"core/core_errors"
	"database/sql"
	"fmt"
	"profiles/data/store"
	"profiles/domain/entities"
	"profiles/domain/values"
)

type SqlDB struct {
	sql *sql.DB
}

func NewSqlDB(sql *sql.DB) (*SqlDB, error) {
	err := createTables(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql}, nil
}

func createTables(sql *sql.DB) error {
	_, err := sql.Exec(`CREATE TABLE IF NOT EXISTS Profile(
		id INTEGER PRIMARYKEY,
		username VARCHAR(255),
		about TEXT,
		avatarPath VARCHAR(255)
	);`)
	if err != nil {
		return fmt.Errorf("while creating profile table: %w", err)
	}
	return nil
}

func (db *SqlDB) CreateProfile(newProfile values.NewProfile) error {
	_, err := db.sql.Exec(`INSERT INTO Profile(id, username, about, avatarPath) values(
		?, ?, ?, ?
	)`, newProfile.Id, newProfile.Username, newProfile.About, newProfile.AvatarPath)
	if err != nil {
		return fmt.Errorf("while inserting into Profile table: %v", err)
	}
	return nil
}

func (db *SqlDB) GetProfile(profileId string) (entities.Profile, error) {
	row := db.sql.QueryRow(`SELECT id, username, about, avatarPath from Profile where id = ?`, profileId)
	profile := entities.Profile{}
	err := row.Scan(&profile.Id, &profile.Username, &profile.About, &profile.AvatarPath)
	if err == sql.ErrNoRows {
		return entities.Profile{}, core_errors.ErrNotFound
	}
	if err != nil {
		return entities.Profile{}, fmt.Errorf("while getting a profile from profile table: %w", err)
	}
	return profile, nil
}

func (db *SqlDB) UpdateProfile(userId string, upd store.DBUpdateData) error {
	_, err := db.sql.Exec(`
	UPDATE Profile SET 
		avatarPath = CASE WHEN ?1 = "" THEN avatarPath ELSE ?1 END,
		about = 	 CASE WHEN ?2 = "" THEN about ELSE ?2 END 
	WHERE id = ?`, upd.AvatarPath, upd.About, userId)
	if err != nil {
		return fmt.Errorf("while updating avatarPath in db: %w", err)
	}
	return nil
}

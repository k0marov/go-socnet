package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/likeable/table_name"
	"github.com/k0marov/socnet/features/profiles/store/models"

	"github.com/k0marov/socnet/features/profiles/store"

	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
)

type SqlDB struct {
	sql       *sql.DB
	TableName table_name.TableName
}

func NewSqlDB(sql *sql.DB) (*SqlDB, error) {
	err := initSQL(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql, TableName: table_name.NewTableName("Profile")}, nil
}

func initSQL(sql *sql.DB) error {
	_, err := sql.Exec(`CREATE TABLE IF NOT EXISTS Profile(
		id INTEGER PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		about TEXT NOT NULL,
		avatarPath VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		return fmt.Errorf("while creating Profile table: %w", err)
	}
	return nil
}

func (db *SqlDB) CreateProfile(newProfile models.ProfileModel) error {
	_, err := db.sql.Exec(`INSERT INTO Profile(id, username, about, avatarPath) values(
		?, ?, ?, ?
	)`, newProfile.Id, newProfile.Username, newProfile.About, newProfile.AvatarPath)
	if err != nil {
		return fmt.Errorf("while inserting into Profile table: %v", err)
	}
	return nil
}

func (db *SqlDB) GetProfile(profileId core_values.UserId) (models.ProfileModel, error) {
	row := db.sql.QueryRow(`
		SELECT id, username, about, avatarPath
		FROM Profile
		WHERE id = ?1`,
		profileId,
	)
	profile := models.ProfileModel{}
	err := row.Scan(&profile.Id, &profile.Username, &profile.About, &profile.AvatarPath)
	if err == sql.ErrNoRows {
		return models.ProfileModel{}, core_errors.ErrNotFound
	}
	if err != nil {
		return models.ProfileModel{}, fmt.Errorf("while getting a profile from profile table: %w", err)
	}

	return profile, nil
}

func (db *SqlDB) UpdateProfile(userId core_values.UserId, upd store.DBUpdateData) error {
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

package sql_db

import (
	"core/core_errors"
	"database/sql"
	"fmt"
	"profiles/domain/entities"
	"profiles/domain/values"
	"profiles/store"
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
		return fmt.Errorf("while creating Profile table: %w", err)
	}
	_, err = sql.Exec(`CREATE TABLE IF NOT EXISTS Follow(
		target_id   INT,
		follower_id INT,
		FOREIGN KEY(target_id) REFERENCES Profile(id) ON DELETE CASCADE,
		FOREIGN KEY(follower_id) REFERENCES Profile(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return fmt.Errorf("while creating Follow table: %w", err)
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

func (db *SqlDB) GetProfile(profileId values.UserId) (entities.Profile, error) {
	row := db.sql.QueryRow(`SELECT id, username, about, avatarPath from Profile where id = ?`, profileId)
	profile := entities.Profile{}
	err := row.Scan(&profile.Id, &profile.Username, &profile.About, &profile.AvatarPath)
	if err == sql.ErrNoRows {
		return entities.Profile{}, core_errors.ErrNotFound
	}
	if err != nil {
		return entities.Profile{}, fmt.Errorf("while getting a profile from profile table: %w", err)
	}

	followsRow := db.sql.QueryRow(`SELECT COUNT(*) from Follow where follower_id = ? `, profileId)
	err = followsRow.Scan(&profile.Follows)
	if err != nil {
		return entities.Profile{}, fmt.Errorf("while getting count of follows: %w", err)
	}
	followersRow := db.sql.QueryRow(`SELECT COUNT(*) from Follow where target_id = ?`, profileId)
	err = followersRow.Scan(&profile.Followers)
	if err != nil {
		return entities.Profile{}, fmt.Errorf("while getting count of followers: %w", err)
	}

	return profile, nil
}

func (db *SqlDB) UpdateProfile(userId values.UserId, upd store.DBUpdateData) error {
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

func (db *SqlDB) GetFollows(userId values.UserId) ([]entities.Profile, error) {
	rows, err := db.sql.Query(`SELECT target_id from Follow where follower_id = ?`, userId)
	if err != nil {
		return []entities.Profile{}, fmt.Errorf("while selecting follows ids from sql: %w", err)
	}
	defer rows.Close()
	result := []entities.Profile{}
	for rows.Next() {
		followedId := ""
		err := rows.Scan(&followedId)
		if err != nil {
			return []entities.Profile{}, fmt.Errorf("error while scanning a row of follows: %w", err)
		}
		followedProfile, err := db.GetProfile(followedId)
		if err != nil {
			return []entities.Profile{}, fmt.Errorf("error while getting the followed profile: %w", err)
		}
		result = append(result, followedProfile)
	}
	return result, nil
}

func (db *SqlDB) IsFollowing(target, follower values.UserId) (bool, error) {
	row := db.sql.QueryRow(`
	SELECT (target_id) FROM Follow WHERE target_id = ? AND follower_id = ? LIMIT 1
	`, target, follower)
	dummyTargetId := ""
	err := row.Scan(&dummyTargetId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("while querying for a Follow: %w", err)
	}
	return true, nil
}

func (db *SqlDB) Follow(target, follower values.UserId) error {
	_, err := db.sql.Exec(`
	INSERT INTO Follow(target_id, follower_id) values (?, ?)
	`, target, follower)
	if err != nil {
		return fmt.Errorf("while inserting a new Follow: %w", err)
	}
	return nil
}

func (db *SqlDB) Unfollow(target, unfollower values.UserId) error {
	_, err := db.sql.Exec(`
	DELETE FROM Follow where target_id = ? AND follower_id = ?
	`, target, unfollower)
	if err != nil {
		return fmt.Errorf("while deleting a Follow: %w", err)
	}
	return nil
}

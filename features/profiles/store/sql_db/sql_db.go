package sql_db

import (
	"database/sql"
	"fmt"

	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"
	"github.com/k0marov/socnet/features/profiles/store"

	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
)

type SqlDB struct {
	sql *sql.DB
}

func NewSqlDB(sql *sql.DB) (*SqlDB, error) {
	err := initSQL(sql)
	if err != nil {
		return nil, err
	}
	return &SqlDB{sql: sql}, nil
}

func initSQL(sql *sql.DB) error {
	_, err := sql.Exec(`CREATE TABLE IF NOT EXISTS Profile(
		id INTEGER PRIMARYKEY,
		username VARCHAR(255) NOT NULL,
		about TEXT NOT NULL,
		avatarPath VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		return fmt.Errorf("while creating Profile table: %w", err)
	}
	_, err = sql.Exec(`CREATE TABLE IF NOT EXISTS Follow(
		target_id   INT NOT NULL,
		follower_id INT NOT NULL,
		FOREIGN KEY(target_id) REFERENCES Profile(id) ON DELETE CASCADE,
		FOREIGN KEY(follower_id) REFERENCES Profile(id) ON DELETE CASCADE
	);`)
	if err != nil {
		return fmt.Errorf("while creating Follow table: %w", err)
	}
	_, err = sql.Exec(`CREATE INDEX IF NOT EXISTS FollowIndex ON Follow(target_id, follower_id)`)
	if err != nil {
		return fmt.Errorf("while creating Follow index: %w", err)
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

func (db *SqlDB) GetProfile(profileId core_values.UserId) (entities.Profile, error) {
	row := db.sql.QueryRow(`
		SELECT id, username, about, avatarPath,
			(SELECT COUNT(*) FROM Follow WHERE follower_id = ?1) AS follows, 
			(SELECT COUNT(*) FROM Follow WHERE target_id = ?1) AS followers 
		FROM Profile
		WHERE id = ?1`,
		profileId,
	)
	profile := entities.Profile{}
	err := row.Scan(&profile.Id, &profile.Username, &profile.About, &profile.AvatarPath, &profile.Follows, &profile.Followers)
	if err == sql.ErrNoRows {
		return entities.Profile{}, core_errors.ErrNotFound
	}
	if err != nil {
		return entities.Profile{}, fmt.Errorf("while getting a profile from profile table: %w", err)
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

func (db *SqlDB) GetFollows(userId core_values.UserId) ([]core_values.UserId, error) {
	rows, err := db.sql.Query(`
		SELECT target_id
		FROM Follow f
		WHERE f.follower_id = ?
	`, userId)
	if err != nil {
		return []core_values.UserId{}, fmt.Errorf("while querying for follows: %w", err)
	}
	defer rows.Close()

	follows := []core_values.UserId{}
	for rows.Next() {
		followId := ""
		err := rows.Scan(&followId)
		if err != nil {
			return []core_values.UserId{}, fmt.Errorf("while scanning for a follow: %w", err)
		}
		follows = append(follows, followId)
	}
	return follows, nil
}

func (db *SqlDB) IsFollowing(target, follower core_values.UserId) (bool, error) {
	row := db.sql.QueryRow(`
	SELECT EXISTS(SELECT 1 FROM Follow WHERE target_id = ? AND follower_id = ?)
	`, target, follower)
	isFollowing := 0
	err := row.Scan(&isFollowing)
	if err != nil {
		return false, fmt.Errorf("while querying for a Follow: %w", err)
	}
	return isFollowing == 1, nil
}

func (db *SqlDB) Follow(target, follower core_values.UserId) error {
	_, err := db.sql.Exec(`
	INSERT INTO Follow(target_id, follower_id) values (?, ?)
	`, target, follower)
	if err != nil {
		return fmt.Errorf("while inserting a new Follow: %w", err)
	}
	return nil
}

func (db *SqlDB) Unfollow(target, unfollower core_values.UserId) error {
	_, err := db.sql.Exec(`
	DELETE FROM Follow where target_id = ? AND follower_id = ?
	`, target, unfollower)
	if err != nil {
		return fmt.Errorf("while deleting a Follow: %w", err)
	}
	return nil
}

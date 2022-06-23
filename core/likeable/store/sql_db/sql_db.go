package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
	"github.com/k0marov/socnet/core/likeable/table_name"
)

type SqlDB struct {
	sql               *sql.DB
	safeLikeableTable string
}

func NewSqlDB(db *sql.DB, targetTable table_name.TableName) (*SqlDB, error) {
	targetName, err := targetTable.Value()
	if err != nil {
		return nil, fmt.Errorf("while getting target table name: %w", err)
	}
	likeableTable := table_name.NewTableName("Likeable" + targetName)
	likeableName, err := likeableTable.Value()
	if err != nil {
		return nil, fmt.Errorf("while generating likeable table name: %w", err)
	}

	err = initSQL(db, targetName, likeableName)
	if err != nil {
		return nil, fmt.Errorf("while initializing sql for likeable %s: %w", targetName, err)
	}
	return &SqlDB{
		sql:               db,
		safeLikeableTable: likeableName,
	}, nil
}

func initSQL(db *sql.DB, verifiedTarget, verifiedLikeable string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + verifiedLikeable + `(
			target_id INT NOT NULL, 
			liker_id INT NOT NULL, 
			FOREIGN KEY(target_id) REFERENCES ` + verifiedTarget + `(id), 
			FOREIGN KEY(liker_Id) REFERENCES Profile(id)
		)
    `)
	if err != nil {
		return fmt.Errorf("while creating table %s: %w", verifiedLikeable, err)
	}
	return nil
}

func (db *SqlDB) IsLiked(target string, liker core_values.UserId) (bool, error) {
	row := db.sql.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM `+db.safeLikeableTable+` WHERE target_id = ? AND liker_id = ?)
	`, target, liker)
	isLiked := 0
	err := row.Scan(&isLiked)
	if err != nil {
		return false, fmt.Errorf("while SELECTing is %s liked: %w", db.safeLikeableTable, err)
	}
	return isLiked == 1, nil
}

func (db *SqlDB) Like(target string, liker core_values.UserId) error {
	_, err := db.sql.Exec(`
		INSERT INTO `+db.safeLikeableTable+`(target_id, liker_id) VALUES(?, ?)
    `, target, liker)
	if err != nil {
		return fmt.Errorf("while INSERTing a new %s: %w", db.safeLikeableTable, err)
	}
	return nil
}

func (db *SqlDB) Unlike(target string, unliker core_values.UserId) error {
	_, err := db.sql.Exec(`
		DELETE FROM `+db.safeLikeableTable+` WHERE target_id = ? AND liker_id = ?
	`, target, unliker)
	if err != nil {
		return fmt.Errorf("while DELETEing a PostLike: %w", err)
	}
	return nil
}

func (db *SqlDB) GetLikesCount(target string) (int, error) {
	row := db.sql.QueryRow(`
		SELECT COUNT(*) FROM `+db.safeLikeableTable+` WHERE target_id = ?
    `, target)
	var likes int
	err := row.Scan(&likes)
	if err != nil {
		return 0, fmt.Errorf("while scanning the likes count: %w", err)
	}
	return likes, nil
}

func (db *SqlDB) GetUserLikesCount(user core_values.UserId) (int, error) {
	row := db.sql.QueryRow(`
		SELECT COUNT(*) FROM `+db.safeLikeableTable+` WHERE liker_id = ?
    `, user)
	var userLikes int
	err := row.Scan(&userLikes)
	if err != nil {
		return 0, fmt.Errorf("while scanning the user likes count: %w", err)
	}
	return userLikes, nil
}

func (db *SqlDB) GetUserLikes(user core_values.UserId) (targetIds []string, err error) {
	rows, err := db.sql.Query(`
		SELECT target_id FROM `+db.safeLikeableTable+` WHERE liker_id = ? 
    `, user)
	if err != nil {
		return []string{}, fmt.Errorf("while SELECTing the target ids that are liked by user: %w", err)
	}
	for rows.Next() {
		var targetId string
		err := rows.Scan(&targetId)
		if err != nil {
			return []string{}, fmt.Errorf("while scanning a target id liked by user: %w", err)
		}
		targetIds = append(targetIds, targetId)
	}
	return targetIds, nil
}
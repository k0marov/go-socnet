package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/socnet/core/core_values"
)

type SqlDB struct {
	sql *sql.DB
}

func NewSqlDB(db *sql.DB, targetTable string) (*SqlDB, error) {
	likeableTable := "Likeable" + targetTable
	err := initSQL(db, likeableTable, targetTable)
	if err != nil {
		return nil, fmt.Errorf("while initializing sql for likeable %s: %w", targetTable, err)
	}
	return &SqlDB{db}, nil
}

func initSQL(db *sql.DB, likeableTable, targetTableName string) error {
	return nil
}

func (db *SqlDB) IsLiked(target string, liker core_values.UserId) (bool, error) {
	panic("unimplemented")
}

func (db *SqlDB) Like(target string, liker core_values.UserId) error {
	panic("unimplemented")
}

func (db *SqlDB) Unlike(target string, liker core_values.UserId) error {
	panic("unimplemented")
}

func (db *SqlDB) GetLikesCount(target string) (int, error) {
	panic("unimplemented")
}

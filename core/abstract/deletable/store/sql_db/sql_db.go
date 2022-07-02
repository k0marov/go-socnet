package sql_db

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
)

type SqlDB struct {
}

func NewSqlDB(db *sql.DB, tableName table_name.TableName) (*SqlDB, error) {
	return &SqlDB{}, nil
}

func (db *SqlDB) Delete(targetId string) error {
	panic("unimplemented")
}

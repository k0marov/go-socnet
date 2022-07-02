package sql_db

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/likeable/table_name"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SqlDB struct {
	sql       *sql.DB
	safeTable string
}

func NewSqlDB(db *sql.DB, targetTable table_name.TableName) (*SqlDB, error) {
	return &SqlDB{}, nil
}

func (db *SqlDB) GetOwner(targetId string) (core_values.UserId, error) {
	panic("unimplemented")
}

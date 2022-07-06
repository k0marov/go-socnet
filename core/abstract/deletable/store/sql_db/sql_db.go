package sql_db

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
)

type SqlDB struct {
	sql             *sql.DB
	safeTargetTable string
}

func NewSqlDB(db *sql.DB, tableName table_name.TableName) (*SqlDB, error) {
	targetTable, err := tableName.Value()
	if err != nil {
		return nil, core_err.Rethrow("getting target table name", err)
	}
	return &SqlDB{sql: db, safeTargetTable: targetTable}, nil
}

func (db *SqlDB) Delete(targetId string) error {
	_, err := db.sql.Exec(`
		DELETE FROM `+db.safeTargetTable+` WHERE id = ?
    `, targetId)
	if err != nil {
		return core_err.Rethrow("deleting a Deletable target", err)
	}
	return nil
}

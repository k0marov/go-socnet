package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
)

type SqlDB struct {
	sql             *sql.DB
	safeTargetTable string
}

func NewSqlDB(db *sql.DB, tableName table_name.TableName) (*SqlDB, error) {
	targetTable, err := tableName.Value()
	if err != nil {
		return nil, fmt.Errorf("while getting target table name: %w", err)
	}
	return &SqlDB{sql: db, safeTargetTable: targetTable}, nil
}

func (db *SqlDB) Delete(targetId string) error {
	_, err := db.sql.Exec(`
		DELETE FROM `+db.safeTargetTable+` WHERE id = ?
    `, targetId)
	if err != nil {
		return fmt.Errorf("while DELETEing from table %s: %w", db.safeTargetTable, err)
	}
	return nil
}

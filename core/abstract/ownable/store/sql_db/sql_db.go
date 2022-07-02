package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SqlDB struct {
	sql             *sql.DB
	safeTargetTable string
}

func NewSqlDB(db *sql.DB, targetTable table_name.TableName) (*SqlDB, error) {
	targetName, err := targetTable.Value()
	if err != nil {
		return nil, fmt.Errorf("while getting target table name: %w", err)
	}

	return &SqlDB{sql: db, safeTargetTable: targetName}, nil
}

func (db *SqlDB) GetOwner(targetId string) (core_values.UserId, error) {
	row := db.sql.QueryRow(`
		SELECT owner_id FROM `+db.safeTargetTable+` WHERE id = ?
    `, targetId)
	var owner core_values.UserId
	err := row.Scan(&owner)
	if err != nil {
		return "", fmt.Errorf("while scanning the owner id: %w", err)
	}
	return owner, nil
}

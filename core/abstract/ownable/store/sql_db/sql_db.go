package sql_db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SqlDB struct {
	sql             *sqlx.DB
	safeTargetTable string
}

func NewSqlDB(db *sqlx.DB, targetTable table_name.TableName) (*SqlDB, error) {
	targetName, err := targetTable.Value()
	if err != nil {
		return nil, core_err.Rethrow("getting target table name", err)
	}

	return &SqlDB{sql: db, safeTargetTable: targetName}, nil
}

func (db *SqlDB) GetOwner(targetId string) (core_values.UserId, error) {
	row := db.sql.QueryRow(`
		SELECT owner_id FROM `+db.safeTargetTable+` WHERE id = ?
    `, targetId)
	var owner core_values.UserId
	err := row.Scan(&owner)
	if err == sql.ErrNoRows {
		return "", core_err.ErrNotFound
	}
	if err != nil {
		return "", core_err.Rethrow("scanning the owner id", err)
	}
	return owner, nil
}

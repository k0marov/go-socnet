package sql_db

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type sqlDB struct {
	sql             *sql.DB
	safeTargetTable string
}

func NewSqlDB(db *sql.DB, targetTable table_name.TableName) (*sqlDB, error) {
	targetName, err := targetTable.Value()
	if err != nil {
		return nil, core_err.Rethrow("getting target table name", err)
	}

	return &sqlDB{sql: db, safeTargetTable: targetName}, nil
}

func (db *sqlDB) GetRecs(user core_values.UserId) ([]string, error) {
	panic("no impl")
}

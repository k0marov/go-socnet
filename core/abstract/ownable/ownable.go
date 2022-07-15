package ownable

import (
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/ownable/service"
	"github.com/k0marov/go-socnet/core/abstract/ownable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
)

type (
	OwnerGetter = service.OwnerGetter
)

type ownable struct {
	GetOwner OwnerGetter
}

func NewOwnable(db *sqlx.DB, tableName table_name.TableName) (ownable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return ownable{}, core_err.Rethrow("opening ownable sql db", err)
	}
	// service
	getOwner := service.NewOwnerGetter(sqlDB.GetOwner)
	return ownable{
		GetOwner: getOwner,
	}, nil
}

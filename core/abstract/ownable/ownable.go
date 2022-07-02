package ownable

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/ownable/service"
	"github.com/k0marov/go-socnet/core/abstract/ownable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
)

type (
	OwnerGetter = service.OwnerGetter
)

type ownable struct {
	GetOwner OwnerGetter
}

func NewOwnable(db *sql.DB, tableName table_name.TableName) (ownable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return ownable{}, fmt.Errorf("while opening ownable sql db: %w", err)
	}
	// service
	getOwner := service.NewOwnerGetter(sqlDB.GetOwner)
	return ownable{
		GetOwner: getOwner,
	}, nil
}

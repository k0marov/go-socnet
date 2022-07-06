package deletable

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/deletable/service"
	"github.com/k0marov/go-socnet/core/abstract/deletable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
)

type (
	Deleter      = service.Deleter
	ForceDeleter = service.ForceDeleter
)

type deletable struct {
	Delete      Deleter
	ForceDelete ForceDeleter
}

func NewDeletable(db *sql.DB, tableName table_name.TableName, ownerGetter ownable.OwnerGetter) (deletable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return deletable{}, core_err.Rethrow("opening sql db for Deletable", err)
	}
	// service
	deleter := service.NewDeleter(ownerGetter, sqlDB.Delete)
	forceDeleter := service.NewForceDeleter(sqlDB.Delete)

	return deletable{
		Delete:      deleter,
		ForceDelete: forceDeleter,
	}, nil
}

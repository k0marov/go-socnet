package deletable

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/deletable/service"
	"github.com/k0marov/go-socnet/core/abstract/deletable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/ownable"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
)

type (
	Deleter = service.Deleter
)

type deletable struct {
	Delete Deleter
}

func NewDeletable(db *sql.DB, tableName table_name.TableName, ownerGetter ownable.OwnerGetter) (deletable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return deletable{}, fmt.Errorf("while opening sql db for Deletable: %w", err)
	}
	// service
	deleter := service.NewDeleter(ownerGetter, sqlDB.Delete)

	return deletable{
		Delete: deleter,
	}, nil
}

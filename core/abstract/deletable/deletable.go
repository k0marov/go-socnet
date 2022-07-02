package deletable

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/deletable/service"
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
	panic("unimplemented")
}

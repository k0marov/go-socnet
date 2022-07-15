package recommendable

import (
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/recommendable/service"
	"github.com/k0marov/go-socnet/core/abstract/recommendable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
)

type (
	RecsGetter  = service.RecsGetter
	RecsUpdater = service.RecsUpdater
)

type Recommendable struct {
	GetRecs    RecsGetter
	UpdateRecs RecsUpdater
}

func NewRecommendable(db *sqlx.DB, tableName table_name.TableName) (Recommendable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return Recommendable{}, core_err.Rethrow("opening Recommendable sql db", err)
	}
	// service
	getRecs := service.NewRecsGetter(sqlDB.GetRecs, sqlDB.GetRandom)
	updateRecs := service.NewRecsUpdater()
	return Recommendable{
		GetRecs:    getRecs,
		UpdateRecs: updateRecs,
	}, nil
}

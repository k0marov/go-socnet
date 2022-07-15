package recommendable

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/recommendable/service"
	"github.com/k0marov/go-socnet/core/abstract/recommendable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
)

type (
	RecsGetter  = service.RecsGetter
	RecsUpdater = service.RecsUpdater
)

type recommendable struct {
	GetRecs    RecsGetter
	UpdateRecs RecsUpdater
}

func NewRecommendable(db *sql.DB, tableName table_name.TableName) (recommendable, error) {
	// store
	sqlDB, err := sql_db.NewSqlDB(db, tableName)
	if err != nil {
		return recommendable{}, core_err.Rethrow("opening recommendable sql db", err)
	}
	// service
	getRecs := service.NewRecsGetter(sqlDB.GetRecs)
	updateRecs := service.NewRecsUpdater()
	return recommendable{
		GetRecs:    getRecs,
		UpdateRecs: updateRecs,
	}, nil
}

package sql_db

import (
	"database/sql"
	"fmt"
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

func initSQL(db *sql.DB, verifiedTarget, verifiedRecommendation string) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + verifiedRecommendation + `(
			recommendation_id INT NOT NULL, 
			user_id INT NOT NULL, 
			FOREIGN KEY(recommendation_id) REFERENCES ` + verifiedTarget + `(id), 
			FOREIGN KEY(user_id) REFERENCES Profile(id)
		)
    `)
	if err != nil {
		return fmt.Errorf("while creating table %s: %w", verifiedRecommendation, err)
	}
	verifiedIndex := verifiedRecommendation + "Index"
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS `+verifiedIndex+` ON `+verifiedRecommendation+` (user_id)
    `, verifiedIndex)
	if err != nil {
		return fmt.Errorf("while creating index %s: %w", verifiedIndex, err)
	}
	return nil
}

func (db *sqlDB) GetRecs(user core_values.UserId, count int) ([]string, error) {
	panic("unimplemented")
}

func (db *sqlDB) StoreRecs(user core_values.UserId, recs []string) error {
	panic("unimplemented")
}

package sql_db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
	"github.com/k0marov/go-socnet/core/helpers"
)

type SqlDB struct {
	sql          *sqlx.DB
	safeRecTable string
}

func NewSqlDB(db *sqlx.DB, targetTable table_name.TableName) (*SqlDB, error) {
	targetName, err := targetTable.Value()
	if err != nil {
		return nil, core_err.Rethrow("getting target table name", err)
	}
	recommendationTable, err := table_name.NewTableName(targetName + "Recommendation").Value()
	if err != nil {
		return nil, core_err.Rethrow("getting recommendation table name", err)
	}
	err = initSQL(db, targetName, recommendationTable)
	if err != nil {
		return nil, core_err.Rethrow("while initializing sql", err)
	}
	return &SqlDB{sql: db, safeRecTable: recommendationTable}, nil
}

func initSQL(db *sqlx.DB, verifiedTarget, verifiedRecommendation string) error {
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

func (db *SqlDB) GetRecs(user core_values.UserId, count int) ([]string, error) {
	var recs []string
	err := db.sql.Select(&recs, `
		SELECT recommendation_id FROM `+db.safeRecTable+` WHERE user_id = ?
		ORDER BY RANDOM()
	    LIMIT ? 
    `, user, count)
	if err != nil {
		return []string{}, core_err.Rethrow("selecting recs from DB", err)
	}
	return recs, nil
}
func (db *SqlDB) GetRandom(count int) ([]string, error) {
	panic("unimplemented")
}

type recModel struct {
	UserId           string `db:"user_id"`
	RecommendationId string `db:"recommendation_id"`
}

func (db *SqlDB) SetRecs(user core_values.UserId, recs []string) error {
	_, err := db.sql.NamedExec(`
		INSERT INTO `+db.safeRecTable+`(recommendation_id, user_id) VALUES (:recommendation_id, :user_id)
    `, helpers.MapForEach(recs, func(rec string) recModel { return recModel{UserId: user, RecommendationId: rec} }))
	if err != nil {
		return core_err.Rethrow("builk inserting recs into DB", err)
	}
	return nil
}

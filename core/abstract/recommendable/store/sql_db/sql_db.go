package sql_db

import (
	"database/sql"
	"fmt"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_err"
	"github.com/k0marov/go-socnet/core/general/core_values"
)

type SqlDB struct {
	sql          *sql.DB
	safeRecTable string
}

func NewSqlDB(db *sql.DB, targetTable table_name.TableName) (*SqlDB, error) {
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

func (db *SqlDB) GetRecs(user core_values.UserId, count int) ([]string, error) {
	panic("unimplemented")
}
func (db *SqlDB) GetRandom(count int) ([]string, error) {
	panic("unimplemented")
}

func (db *SqlDB) SetRecs(user core_values.UserId, recs []string) error {
	panic("unimplemented")
	//statement := "INSERT INTO " + db.safeRecTable + "(recommendation_id, user_id) VALUES " + strings.Repeat("(?, ?),", len(recs))
	//statement = statement[:len(statement)-2] // remove the trailing comma
	//statement += ";"
	//var values []any
	//for _, rec := range recs {
	//	values = append(values, rec)
	//	values = append(values, user)
	//}
	//db.sql.Exec(statement, values...)
	//return nil
}

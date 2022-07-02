package sql_db_test

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/deletable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

var targetTblName = table_name.NewTableName("Target")

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	db.Close() // this will make all calls to db throw
	t.Run("Delete", func(t *testing.T) {
		err := sqlDB.Delete(RandomId())
		AssertSomeError(t, err)
	})
}

func TestSqlDB_Injection(t *testing.T) {
	db := OpenSqliteDB(t)
	_, err := sql_db.NewSqlDB(db, table_name.NewTableName("'; DROP TABLE Students; --"))
	AssertSomeError(t, err)
}

func TestSqlDB(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)

	targetTable, err := targetTblName.Value()
	AssertNoError(t, err)

	targetId := createTargetEntity(t, db)
	err = sqlDB.Delete(targetId)
	AssertNoError(t, err)

	row := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM `+targetTable+` WHERE id = ?); 
    `, targetId)
	var stillExists bool
	row.Scan(&stillExists)

	Assert(t, stillExists, false, "target still exists in db")
}

func setupSqlDB(t testing.TB, db *sql.DB) *sql_db.SqlDB {
	t.Helper()
	targetTable, err := targetTblName.Value()
	AssertNoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + targetTable + `(
		    id INTEGER PRIMARY KEY 
		)
    `)
	AssertNoError(t, err)
	sqlDB, err := sql_db.NewSqlDB(db, targetTblName)
	AssertNoError(t, err)
	return sqlDB
}

func createTargetEntity(t testing.TB, db *sql.DB) (id string) {
	t.Helper()
	targetTable, err := targetTblName.Value()
	AssertNoError(t, err)
	id = RandomId()
	_, err = db.Exec(`
		INSERT INTO `+targetTable+`(id) VALUES (?)
    `, id)
	AssertNoError(t, err)
	return
}

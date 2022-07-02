package sql_db_test

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/ownable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	"github.com/k0marov/go-socnet/core/general/core_values"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"github.com/k0marov/go-socnet/features/profiles/domain/models"
	profiles_db "github.com/k0marov/go-socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

var targetTblName = table_name.NewTableName("Target")

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	db.Close() // this will make all calls to db throw
	t.Run("GetOwner", func(t *testing.T) {
		_, err := sqlDB.GetOwner(RandomId())
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
	profilesDB, _ := profiles_db.NewSqlDB(db)

	ownerId := RandomId()

	profilesDB.CreateProfile(models.ProfileModel{
		Id: ownerId,
	})
	targetId := createTargetEntity(t, db, ownerId)

	gotOwner, err := sqlDB.GetOwner(targetId)
	AssertNoError(t, err)
	Assert(t, ownerId, gotOwner, "returned owner")

}

func setupSqlDB(t testing.TB, db *sql.DB) *sql_db.SqlDB {
	t.Helper()
	targetTable, err := targetTblName.Value()
	AssertNoError(t, err)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + targetTable + `(
		    id INTEGER PRIMARY KEY, 
			owner_id INT NOT NULL, 
			FOREIGN KEY(owner_id) REFERENCES Profile(id)
		)
    `)
	AssertNoError(t, err)
	sqlDB, err := sql_db.NewSqlDB(db, targetTblName)
	AssertNoError(t, err)
	return sqlDB
}

func createTargetEntity(t testing.TB, db *sql.DB, owner core_values.UserId) (id string) {
	t.Helper()
	targetTable, err := targetTblName.Value()
	AssertNoError(t, err)
	id = RandomId()
	_, err = db.Exec(`
		INSERT INTO `+targetTable+`(id, owner_id) VALUES (?, ?)
    `, id, owner)
	AssertNoError(t, err)
	return
}

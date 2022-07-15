package sql_db_test

import (
	"github.com/jmoiron/sqlx"
	"github.com/k0marov/go-socnet/core/abstract/recommendable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	profiles_db "github.com/k0marov/go-socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"sort"
	"testing"
)

var targetTblName = table_name.NewTableName("Target")

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	db.Close() // this will make all calls to db throw
	t.Run("GetRecs", func(t *testing.T) {
		_, err := sqlDB.GetRecs(RandomId(), RandomInt())
		AssertSomeError(t, err)
	})
	t.Run("GetRandom", func(t *testing.T) {
		_, err := sqlDB.GetRandom(RandomInt())
		AssertSomeError(t, err)
	})
	t.Run("SetRecs", func(t *testing.T) {
		err := sqlDB.SetRecs(RandomId(), []string{RandomId(), RandomId()})
		AssertSomeError(t, err)
	})
}
func TestSqlDB(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	profilesDB, err := profiles_db.NewSqlDB(db)
	AssertNoError(t, err)

	profile1 := RandomProfileModel()
	profilesDB.CreateProfile(profile1)
	profile2 := RandomProfileModel()
	profilesDB.CreateProfile(profile2)

	// add a hundred recommendations for the first profile
	var targets []string
	for i := 0; i < 100; i++ {
		targets = append(targets, createTargetEntity(t, db))
	}
	sort.Strings(targets)
	err = sqlDB.SetRecs(profile1.Id, targets)
	AssertNoError(t, err)

	gotRecs, err := sqlDB.GetRecs(profile1.Id, 100)
	AssertNoError(t, err)
	sort.Strings(gotRecs)
	Assert(t, gotRecs, targets, "returned recommendations")

	// assert randomness
	gotRecs2, err := sqlDB.GetRecs(profile1.Id, 100)
	Assert(t, reflect.DeepEqual(gotRecs, gotRecs2), false, "the returned values are not equal")

	// assert that limiting count works
	gotRecsLimited, err := sqlDB.GetRecs(profile1.Id, 10)
	AssertNoError(t, err)
	Assert(t, len(gotRecsLimited), 10, "length of limited recs")
}

func TestSqlDB_Injection(t *testing.T) {
	db := OpenSqliteDB(t)
	_, err := sql_db.NewSqlDB(db, table_name.NewTableName("'; DROP TABLE Students; --"))
	AssertSomeError(t, err)
}

func setupSqlDB(t testing.TB, db *sqlx.DB) *sql_db.SqlDB {
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

func createTargetEntity(t testing.TB, db *sqlx.DB) (id string) {
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

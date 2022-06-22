package sql_db_test

import (
	"database/sql"
	"github.com/k0marov/socnet/core/likeable/store/sql_db"
	. "github.com/k0marov/socnet/core/test_helpers"
	profiles_db "github.com/k0marov/socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	db.Close() // this will make all calls to db throw
	t.Run("IsLiked", func(t *testing.T) {
		_, err := sqlDB.IsLiked(RandomId(), RandomId())
		AssertSomeError(t, err)
	})
	t.Run("Like", func(t *testing.T) {
		err := sqlDB.Like(RandomId(), RandomId())
		AssertSomeError(t, err)
	})
	t.Run("Unlike", func(t *testing.T) {
		err := sqlDB.Unlike(RandomId(), RandomId())
		AssertSomeError(t, err)
	})
	t.Run("GetLikesCount", func(t *testing.T) {
		_, err := sqlDB.GetLikesCount(RandomId())
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB := setupSqlDB(t, db)
	profilesDB, err := profiles_db.NewSqlDB(db)
	AssertNoError(t, err)

	t.Run("liking/unliking", func(t *testing.T) {
		// create a target entity
		targetId := createTargetEntity(t, db)
		// create a profile
		profile := RandomNewProfile()
		profilesDB.CreateProfile(profile)

		// define a helper
		assertLikedValue := func(t testing.TB, value bool) {
			got, err := sqlDB.IsLiked(targetId, profile.Id)
			AssertNoError(t, err)
			Assert(t, got, value, "the isLiked value")
		}

		// assert target is not liked from profile
		assertLikedValue(t, false)
		// like it
		err := sqlDB.Like(targetId, profile.Id)
		AssertNoError(t, err)
		// assert it is liked
		assertLikedValue(t, true)
		// unlike it
		err = sqlDB.Unlike(targetId, profile.Id)
		AssertNoError(t, err)
		// assert it is not liked
		assertLikedValue(t, false)

	})
	t.Run("liking a non-existing target", func(t *testing.T) {

	})
}

func setupSqlDB(t testing.TB, db *sql.DB) *sql_db.SqlDB {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Target(
		    id INTEGER PRIMARY KEY
		)
    `)
	AssertNoError(t, err)
	sqlDB, err := sql_db.NewSqlDB(db, "Target")
	AssertNoError(t, err)
	return sqlDB
}

func createTargetEntity(t testing.TB, db *sql.DB) (id string) {
	t.Helper()
	id = RandomId()
	_, err := db.Exec(`
		INSERT INTO Target(id) VALUES (?)
    `, id)
	AssertNoError(t, err)
	return
}

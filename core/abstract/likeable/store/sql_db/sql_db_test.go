package sql_db_test

import (
	"database/sql"
	"github.com/k0marov/go-socnet/core/abstract/likeable/store/sql_db"
	"github.com/k0marov/go-socnet/core/abstract/table_name"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"

	profiles_db "github.com/k0marov/go-socnet/features/profiles/store/sql_db"
	_ "github.com/mattn/go-sqlite3"
)

var targetTblName = table_name.NewTableName("Target")

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
	t.Run("GetUserLikesCount", func(t *testing.T) {
		_, err := sqlDB.GetUserLikesCount(RandomId())
		AssertSomeError(t, err)
	})
	t.Run("GetUserLikes", func(t *testing.T) {
		_, err := sqlDB.GetUserLikes(RandomId())
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
	profilesDB, err := profiles_db.NewSqlDB(db)
	AssertNoError(t, err)

	t.Run("liking/unliking", func(t *testing.T) {
		// create a target entity
		targetId := createTargetEntity(t, db)
		// create a profile
		profile := RandomProfileModel()
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
	t.Run("liking 1 target from many profiles", func(t *testing.T) {
		targetId := createTargetEntity(t, db)
		const count = 100
		for i := 0; i < count; i++ {
			profile := RandomProfileModel()
			profilesDB.CreateProfile(profile)
			err := sqlDB.Like(targetId, profile.Id)
			AssertNoError(t, err)

			likes, err := sqlDB.GetLikesCount(targetId)
			AssertNoError(t, err)
			Assert(t, likes, i+1, "number of likes")
		}
	})
	t.Run("liking many targets from 1 profile", func(t *testing.T) {
		const count = 100
		profile := RandomProfileModel()
		profilesDB.CreateProfile(profile)

		var targets []string
		for i := 0; i < count; i++ {
			target := createTargetEntity(t, db)
			targets = append(targets, target)
			err := sqlDB.Like(target, profile.Id)
			AssertNoError(t, err)

			userLikesCount, err := sqlDB.GetUserLikesCount(profile.Id)
			AssertNoError(t, err)
			Assert(t, userLikesCount, i+1, "number of targets liked by user")

			userLikes, err := sqlDB.GetUserLikes(profile.Id)
			AssertNoError(t, err)
			Assert(t, userLikes, targets, "targets liked by user")
		}
	})
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

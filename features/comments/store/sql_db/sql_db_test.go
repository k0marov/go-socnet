package sql_db_test

import (
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/comments/store/sql_db"
	"testing"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	db := OpenSqliteDB(t)
	sqlDB, err := sql_db.NewSqlDB(db)
	AssertNoError(t, err)
	db.Close() // this will make all calls to db throw
	t.Run("IsLiked", func(t *testing.T) {
		_, err := sqlDB.IsLiked(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Like", func(t *testing.T) {
		err := sqlDB.Like(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Unlike", func(t *testing.T) {
		err := sqlDB.Unlike(RandomString(), RandomString())
		AssertSomeError(t, err)
	})
	t.Run("GetComments", func(t *testing.T) {
		_, err := sqlDB.GetComments(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("Create", func(t *testing.T) {
		_, err := sqlDB.Create(RandomNewComment())
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {

}

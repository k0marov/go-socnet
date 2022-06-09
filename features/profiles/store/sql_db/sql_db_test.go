package sql_db_test

import (
	"core/core_errors"
	. "core/test_helpers"
	"database/sql"
	"profiles/data/store"
	"profiles/data/store/sql_db"
	"profiles/domain/entities"
	"profiles/domain/values"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	sql := OpenSqliteDB(t)
	sut, _ := sql_db.NewSqlDB(sql)
	sql.Close() // this will force all calls to throw errors
	t.Run("GetProfile", func(t *testing.T) {
		_, err := sut.GetProfile(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("CreateProfile", func(t *testing.T) {
		err := sut.CreateProfile(values.NewProfile{Profile: RandomProfile()})
		AssertSomeError(t, err)
	})
	t.Run("UpdateProfile", func(t *testing.T) {
		err := sut.UpdateProfile(RandomString(), store.DBUpdateData{About: RandomString(), AvatarPath: RandomString()})
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	t.Run("creating and reading profiles", func(t *testing.T) {
		profileCount := 10

		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		// create 10 random profiles
		profiles := []entities.Profile{}
		for i := 0; i < profileCount; i++ {
			profiles = append(profiles, RandomProfile())
		}

		// add them to db
		for _, profile := range profiles {
			db.CreateProfile(values.NewProfile{Profile: profile})
		}

		// assert they can be found in the database
		for _, profile := range profiles {
			gotProfile, err := db.GetProfile(profile.Id)
			AssertNoError(t, err)
			Assert(t, gotProfile, profile, "profile stored in db")
		}

		// assert querying for unexisting profile returns ErrNotFound
		_, err = db.GetProfile(strconv.Itoa(9999))
		AssertError(t, err, core_errors.ErrNotFound)
	})
	t.Run("updating profile", func(t *testing.T) {
		profile1 := RandomProfile()
		profile2 := RandomProfile()

		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		// insert both profiles into database
		db.CreateProfile(values.NewProfile{Profile: profile1})
		db.CreateProfile(values.NewProfile{Profile: profile2})

		// update first profile
		newAvatar := RandomString()
		newAbout := RandomString()
		err = db.UpdateProfile(profile1.Id, store.DBUpdateData{AvatarPath: newAvatar, About: newAbout})
		AssertNoError(t, err)

		// update second profile with empty values (it shouldn't be updated)
		err = db.UpdateProfile(profile2.Id, store.DBUpdateData{})
		AssertNoError(t, err)

		// get the updated profile, it should be changed
		wantUpdatedProfile1 := entities.Profile{
			Id:         profile1.Id,
			Username:   profile1.Username,
			About:      newAbout,
			AvatarPath: newAvatar,
		}
		updatedProfile1, err := db.GetProfile(profile1.Id)
		AssertNoError(t, err)
		Assert(t, updatedProfile1, wantUpdatedProfile1, "the updated profile")

		// get the unaffected profile, it shouldn't be changed
		gotProfile2, err := db.GetProfile(profile2.Id)
		AssertNoError(t, err)
		Assert(t, gotProfile2, profile2, "the unaffected profile")
	})
}

func OpenSqliteDB(t testing.TB) *sql.DB {
	t.Helper()
	sql, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("error while opening in-memory database: %v", err)
	}
	return sql
}

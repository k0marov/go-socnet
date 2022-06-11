package sql_db_test

import (
	"core/core_errors"
	. "core/test_helpers"
	"database/sql"
	"profiles/domain/entities"
	"profiles/domain/values"
	"profiles/store"
	"profiles/store/sql_db"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlDB_ErrorHandling(t *testing.T) {
	sql := OpenSqliteDB(t)
	sut, err := sql_db.NewSqlDB(sql)
	AssertNoError(t, err)
	sql.Close() // this will force all calls to throw errors
	t.Run("GetProfile", func(t *testing.T) {
		_, err := sut.GetProfile(RandomString())
		AssertSomeError(t, err)
	})
	t.Run("CreateProfile", func(t *testing.T) {
		err := sut.CreateProfile(RandomProfile())
		AssertSomeError(t, err)
	})
	t.Run("UpdateProfile", func(t *testing.T) {
		err := sut.UpdateProfile(RandomString(), store.DBUpdateData{About: RandomString(), AvatarPath: RandomString()})
		AssertSomeError(t, err)
	})
	t.Run("IsFollowing", func(t *testing.T) {
		_, err := sut.IsFollowing("42", "33")
		AssertSomeError(t, err)
	})
	t.Run("Follow", func(t *testing.T) {
		err := sut.Follow("42", "33")
		AssertSomeError(t, err)
	})
	t.Run("Unfollow", func(t *testing.T) {
		err := sut.Unfollow("42", "33")
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
			db.CreateProfile(profile)
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
		db.CreateProfile(profile1)
		db.CreateProfile(profile2)

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
			Follows:    profile1.Follows,
			Followers:  profile1.Followers,
		}
		updatedProfile1, err := db.GetProfile(profile1.Id)
		AssertNoError(t, err)
		Assert(t, updatedProfile1, wantUpdatedProfile1, "the updated profile")

		// get the unaffected profile, it shouldn't be changed
		gotProfile2, err := db.GetProfile(profile2.Id)
		AssertNoError(t, err)
		Assert(t, gotProfile2, profile2, "the unaffected profile")
	})

	t.Run("IsFollowing(), Follow(), Unfollow()", func(t *testing.T) {
		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		profile1 := RandomProfile()
		profile2 := RandomProfile()

		// create 2 profiles
		db.CreateProfile(profile1)
		db.CreateProfile(profile2)

		// they are not following each other
		assertFollows := func(target, follower values.UserId, value bool) {
			follows, err := db.IsFollowing(target, follower)
			AssertNoError(t, err)
			Assert(t, follows, value, "returned value")
		}
		assertFollows(profile1.Id, profile2.Id, false)
		assertFollows(profile2.Id, profile1.Id, false)

		// make 1-st profile follow the 2-nd profile
		err = db.Follow(profile2.Id, profile1.Id)
		AssertNoError(t, err)

		// now 1-st profile should follow the 2-nd profile
		assertFollows(profile2.Id, profile1.Id, true)
		// and 2-nd profile should still not follow the 1-st profile
		assertFollows(profile1.Id, profile2.Id, false)

		// now call unfollow
		err = db.Unfollow(profile2.Id, profile1.Id)
		AssertNoError(t, err)

		// now profile1 shouldn't follow profile2
		assertFollows(profile1.Id, profile2.Id, false)
		assertFollows(profile2.Id, profile1.Id, false)

	})
	t.Run("following many profiles: Follows(), GetProfile()", func(t *testing.T) {
		// followingProfile := RandomProfile()

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

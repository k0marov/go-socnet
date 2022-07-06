package sql_db_test

import (
	"github.com/k0marov/go-socnet/core/general/core_err"
	. "github.com/k0marov/go-socnet/core/helpers/test_helpers"
	"testing"

	"github.com/k0marov/go-socnet/features/profiles/domain/models"
	"github.com/k0marov/go-socnet/features/profiles/store"
	"github.com/k0marov/go-socnet/features/profiles/store/sql_db"

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
		err := sut.CreateProfile(RandomProfileModel())
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
		profiles := []models.ProfileModel{}
		for i := 0; i < profileCount; i++ {
			profiles = append(profiles, RandomProfileModel())
		}

		// add them to db
		for _, profile := range profiles {
			db.CreateProfile(profile)
		}

		// assert they can be found in the database
		for _, profile := range profiles {
			wantProfile := models.ProfileModel{
				Id:         profile.Id,
				Username:   profile.Username,
				About:      profile.About,
				AvatarPath: profile.AvatarPath,
			}
			gotProfile, err := db.GetProfile(profile.Id)
			AssertNoError(t, err)
			Assert(t, gotProfile, wantProfile, "profile stored in db")
		}

		// assert querying for unexisting profile returns ErrNotFound
		_, err = db.GetProfile("9999")
		AssertError(t, err, core_err.ErrNotFound)
	})
	t.Run("updating profile", func(t *testing.T) {
		newProfile1 := RandomProfileModel()
		newProfile2 := RandomProfileModel()

		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		// insert both profiles into database
		db.CreateProfile(newProfile1)
		db.CreateProfile(newProfile2)

		// update first profile
		newAvatar := RandomString()
		newAbout := RandomString()
		err = db.UpdateProfile(newProfile1.Id, store.DBUpdateData{AvatarPath: newAvatar, About: newAbout})
		AssertNoError(t, err)

		// update second profile with empty values (it shouldn't be updated)
		err = db.UpdateProfile(newProfile2.Id, store.DBUpdateData{})
		AssertNoError(t, err)

		// get the updated profile, it should be changed
		wantUpdatedProfile1 := models.ProfileModel{
			Id:         newProfile1.Id,
			Username:   newProfile1.Username,
			About:      newAbout,
			AvatarPath: newAvatar,
		}
		updatedProfile1, err := db.GetProfile(newProfile1.Id)
		AssertNoError(t, err)
		Assert(t, updatedProfile1, wantUpdatedProfile1, "the updated profile")

		// get the unaffected profile, it shouldn't be changed
		gotProfile2, err := db.GetProfile(newProfile2.Id)
		AssertNoError(t, err)
		Assert(t, gotProfile2, newProfile2, "the unaffected profile")
	})
}

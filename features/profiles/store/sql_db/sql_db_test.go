package sql_db_test

import (
	"github.com/k0marov/socnet/core/core_errors"
	"github.com/k0marov/socnet/core/core_values"
	. "github.com/k0marov/socnet/core/test_helpers"
	"github.com/k0marov/socnet/features/profiles/domain/entities"
	"github.com/k0marov/socnet/features/profiles/domain/values"
	"github.com/k0marov/socnet/features/profiles/store"
	"github.com/k0marov/socnet/features/profiles/store/sql_db"
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
		err := sut.CreateProfile(RandomNewProfile())
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
	t.Run("GetFollows", func(t *testing.T) {
		_, err := sut.GetFollows("42")
		AssertSomeError(t, err)
	})
}

func TestSqlDB(t *testing.T) {
	t.Run("creating and reading profiles", func(t *testing.T) {
		profileCount := 10

		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		// create 10 random profiles
		profiles := []values.NewProfile{}
		for i := 0; i < profileCount; i++ {
			profiles = append(profiles, RandomNewProfile())
		}

		// add them to db
		for _, profile := range profiles {
			db.CreateProfile(profile)
		}

		// assert they can be found in the database
		for _, profile := range profiles {
			gotProfile, err := db.GetProfile(profile.Id)
			wantProfile := entities.Profile{
				Id:         profile.Id,
				Username:   profile.Username,
				About:      profile.About,
				AvatarPath: profile.AvatarPath,
				Followers:  0,
				Follows:    0,
			}
			AssertNoError(t, err)
			Assert(t, gotProfile, wantProfile, "profile stored in db")
		}

		// assert querying for unexisting profile returns ErrNotFound
		_, err = db.GetProfile("9999")
		AssertError(t, err, core_errors.ErrNotFound)
	})
	t.Run("updating profile", func(t *testing.T) {
		newProfile1 := RandomNewProfile()
		newProfile2 := RandomNewProfile()
		profile1 := entities.Profile{
			Id:         newProfile1.Id,
			Username:   newProfile1.Username,
			About:      newProfile1.About,
			AvatarPath: newProfile1.AvatarPath,
		}
		profile2 := entities.Profile{
			Id:         newProfile2.Id,
			Username:   newProfile2.Username,
			About:      newProfile2.About,
			AvatarPath: newProfile2.AvatarPath,
		}

		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		// insert both profiles into database
		db.CreateProfile(newProfile1)
		db.CreateProfile(newProfile2)

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
	t.Run("following profiles", func(t *testing.T) {
		db, err := sql_db.NewSqlDB(OpenSqliteDB(t))
		AssertNoError(t, err)

		profile1 := RandomNewProfile()
		profile2 := RandomNewProfile()

		// create 2 profiles
		db.CreateProfile(profile1)
		db.CreateProfile(profile2)

		assertFollows := func(t testing.TB, target, follower core_values.UserId, shouldFollow bool) {
			t.Helper()
			isFollowing, err := db.IsFollowing(target, follower)
			AssertNoError(t, err)
			Assert(t, isFollowing, shouldFollow, "returned value of IsFollowing")

			follows, err := db.GetFollows(follower)
			AssertNoError(t, err)
			if shouldFollow {
				AssertFatal(t, len(follows), 1, "number of profiles the 'follower' follows")
				Assert(t, follows[0], target, "the id of followed profile")
			} else {
				AssertFatal(t, len(follows), 0, "number of profiles the 'follower' follows")
			}

			followerProfile, err := db.GetProfile(follower)
			AssertNoError(t, err)
			if shouldFollow {
				Assert(t, followerProfile.Follows, 1, "amount of profiles the 'follower' follows")
			} else {
				Assert(t, followerProfile.Follows, 0, "amount of profiles the 'follower' follows")
			}

			targetProfile, err := db.GetProfile(target)
			AssertNoError(t, err)
			if shouldFollow {
				Assert(t, targetProfile.Followers, 1, "amount of followers on the 'target' profile")
			} else {
				Assert(t, targetProfile.Followers, 0, "amount of followers on the 'target' profile")
			}

		}
		// they are not following each other
		assertFollows(t, profile1.Id, profile2.Id, false)
		assertFollows(t, profile2.Id, profile1.Id, false)

		// make 1-st profile follow the 2-nd profile
		err = db.Follow(profile2.Id, profile1.Id)
		AssertNoError(t, err)

		// now 1-st profile should follow the 2-nd profile
		assertFollows(t, profile2.Id, profile1.Id, true)
		// and 2-nd profile should still not follow the 1-st profile
		assertFollows(t, profile1.Id, profile2.Id, false)

		// call unfollow
		err = db.Unfollow(profile2.Id, profile1.Id)
		AssertNoError(t, err)

		// now profile1 shouldn't follow profile2
		assertFollows(t, profile1.Id, profile2.Id, false)
		assertFollows(t, profile2.Id, profile1.Id, false)

	})
	t.Run("following many profiles", func(t *testing.T) {
		db, _ := sql_db.NewSqlDB(OpenSqliteDB(t))

		mainProfileNew := RandomNewProfile()
		db.CreateProfile(mainProfileNew)

		profileCount := 100
		// create 10 random otherProfiles
		otherProfiles := []values.NewProfile{}
		for i := 0; i < profileCount; i++ {
			newProfile := RandomNewProfile()
			otherProfiles = append(otherProfiles, newProfile)
			db.CreateProfile(newProfile)
		}

		// follow them from the main profile
		for i, otherProfile := range otherProfiles {
			db.Follow(otherProfile.Id, mainProfileNew.Id)
			currentMainProfile, _ := db.GetProfile(mainProfileNew.Id)
			Assert(t, currentMainProfile.Follows, i+1, "number of profiles that main follows")
			isFollowing, _ := db.IsFollowing(otherProfile.Id, mainProfileNew.Id)
			Assert(t, isFollowing, true, "other profile is followed")
		}

		// make them followers of the main profile
		for i, otherProfile := range otherProfiles {
			db.Follow(mainProfileNew.Id, otherProfile.Id)
			currentMainProfile, _ := db.GetProfile(mainProfileNew.Id)
			Assert(t, currentMainProfile.Followers, i+1, "number of followers of the main profile")
			isFollowing, _ := db.IsFollowing(mainProfileNew.Id, otherProfile.Id)
			Assert(t, isFollowing, true, "other profile is a follower")
		}
	})
}

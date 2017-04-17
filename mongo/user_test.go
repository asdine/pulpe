package mongo_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func newUserID() string {
	return bson.NewObjectId().Hex()
}

// Ensure users can be created and retrieved.
func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("New", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()
		s := session.UserService()

		u := pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
		}

		// Create new user.
		user, err := s.CreateUser(&u)
		require.NoError(t, err)
		require.NotZero(t, user.ID)
		require.Equal(t, "jonsnow", user.Login)

		// Retrieve user and compare.
		other, err := s.User(user.ID)
		require.NoError(t, err)
		require.Equal(t, user, other)
	})

	t.Run("Login conflict", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()
		s := session.UserService()

		u := pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
		}

		// Create new user.
		user, err := s.CreateUser(&u)
		require.NoError(t, err)
		require.Equal(t, user.Login, "jonsnow")

		// Create second user with the same login.
		u = pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon-snow@wall.com",
		}

		user, err = s.CreateUser(&u)
		require.NoError(t, err)
		require.Equal(t, "jonsnow1", user.Login)
	})

	t.Run("Email conflict", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()
		s := session.UserService()

		u := pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
		}

		// Create new user.
		_, err := s.CreateUser(&u)
		require.NoError(t, err)

		// Create second user with the same email.
		u = pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
		}

		_, err = s.CreateUser(&u)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrEmailConflict, err)
	})
}

// Ensure users can be retrieved.
func TestUserService_User(t *testing.T) {
	t.Parallel()
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.UserService()

	t.Run("OK", func(t *testing.T) {
		u := pulpe.UserCreation{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
		}

		// Create new user.
		user, err := s.CreateUser(&u)
		require.NoError(t, err)

		// Retrieve user and compare.
		other, err := s.User(user.ID)
		require.NoError(t, err)
		require.Equal(t, user, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to fetch a user that doesn't exist.
		_, err := s.User("something")
		require.Equal(t, pulpe.ErrUserNotFound, err)
	})
}

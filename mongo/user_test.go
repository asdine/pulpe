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
func TestUserService_Register(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user.
		user, err := s.Register(&u)
		require.NoError(t, err)
		require.NotZero(t, user.ID)
		require.Equal(t, "jonsnow", user.Login)

		// Retrieve user and compare.
		other, err := s.User(user.ID)
		require.NoError(t, err)
		require.Equal(t, user, other)
	})

	t.Run("Login conflict", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user.
		user, err := s.Register(&u)
		require.NoError(t, err)
		require.Equal(t, user.Login, "jonsnow")

		// Create second user with the same login.
		u = pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon-snow@wall.com",
			Password: "ygritte",
		}

		user, err = s.Register(&u)
		require.NoError(t, err)
		require.Equal(t, "jonsnow1", user.Login)
	})

	t.Run("Email conflict", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user.
		_, err := s.Register(&u)
		require.NoError(t, err)

		// Create second user with the same email.
		u = pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		_, err = s.Register(&u)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserEmailConflict, err)
	})
}

// Ensure users can be retrieved.
func TestUserService_User(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user
		user, err := s.Register(&u)
		require.NoError(t, err)

		// Retrieve user and compare
		other, err := s.User(user.ID)
		require.NoError(t, err)
		require.Equal(t, user, other)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.UserService()

		// Trying to fetch a user that doesn't exist
		_, err := s.User("something")
		require.Equal(t, pulpe.ErrUserNotFound, err)
	})
}

func TestUserSessionService_Login(t *testing.T) {
	t.Run("WithEmailOK", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user.
		user, err := s.Register(&u)
		require.NoError(t, err)

		// Login user.
		userSession, err := sessions.NoAuth.UserSessionService().Login(u.Email, u.Password)
		require.NoError(t, err)
		require.Equal(t, userSession.UserID, user.ID)
	})

	t.Run("WithLoginOK", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		// Create new user.
		user, err := s.Register(&u)
		require.NoError(t, err)

		// Login user.
		userSession, err := sessions.NoAuth.UserSessionService().Login(user.Login, u.Password)
		require.NoError(t, err)
		require.Equal(t, userSession.UserID, user.ID)
	})

	t.Run("WithBadEmail", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserSessionService()

		_, err := s.Login("someone@email.com", "passwd")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("WithBadLogin", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserSessionService()

		_, err := s.Login("someone", "passwd")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("WithBadPassword", func(t *testing.T) {
		sessions, cleanup := MustGetSessions(t)
		defer cleanup()

		s := sessions.NoAuth.UserService()

		u := pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "ygritte",
		}

		user, err := s.Register(&u)
		require.NoError(t, err)

		_, err = sessions.NoAuth.UserSessionService().Login(user.Login, "passwd")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})
}

func TestUserSessionService_CreateSession(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.NoAuth.UserSessionService()

		u := pulpe.User{
			ID: "id",
		}

		us, err := s.CreateSession(&u)
		require.NoError(t, err)
		require.Equal(t, u.ID, us.UserID)
		require.True(t, us.UpdatedAt.Before(us.ExpiresAt))
		require.NotNil(t, us.ID)
	})
}

func TestUserSessionService_GetSession(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.NoAuth.UserSessionService()

		u := pulpe.User{
			ID: "id",
		}

		us, err := s.CreateSession(&u)
		require.NoError(t, err)

		sess, err := s.GetSession(us.ID)
		require.NoError(t, err)
		require.Equal(t, us, sess)
	})

	t.Run("UnknownSession", func(t *testing.T) {
		s := sessions.NoAuth.UserSessionService()

		_, err := s.GetSession("somesid")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserSessionUnknownID, err)
	})
}

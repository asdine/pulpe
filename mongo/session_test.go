package mongo_test

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/mongo"
	"github.com/stretchr/testify/require"
)

const userPassword = "passw0rd"

type Session struct {
	*mongo.Session
}

type Sessions struct {
	NoAuth *Session
	Red    *Session
	Blue   *Session
	Green  *Session
}

func (s *Sessions) Close() {
	s.NoAuth.Close()
	s.Red.Close()
	s.Blue.Close()
	s.Green.Close()
}

var sessions *Sessions

func MustGetSessions(t require.TestingT) (*Sessions, func()) {
	if sessions == nil {
		sessions = &Sessions{
			NoAuth: &Session{Session: client.Connect()},
			Red:    getSessionAs(t, "Red"),
			Blue:   getSessionAs(t, "Blue"),
			Green:  getSessionAs(t, "Green"),
		}
	}

	return sessions, func() {
		_, err := client.Session.DB("").C("users").RemoveAll(bson.M{
			"fullName": bson.M{
				"$nin": []string{"Red", "Blue", "Green"},
			},
		})
		require.NoError(t, err)

		_, err = client.Session.DB("").C("boards").RemoveAll(nil)
		require.NoError(t, err)

		_, err = client.Session.DB("").C("lists").RemoveAll(nil)
		require.NoError(t, err)

		_, err = client.Session.DB("").C("cards").RemoveAll(nil)
		require.NoError(t, err)
	}
}

func getSessionAs(t require.TestingT, name string) *Session {
	session := client.Connect()
	user := createUser(t, session, name)
	us, err := session.UserSessionService().CreateSession(user)
	require.NoError(t, err)
	session.SetAuthToken(us.ID)
	return &Session{Session: session}
}

func createUser(t require.TestingT, s *mongo.Session, name string) *pulpe.User {
	user, err := s.UserService().Register(&pulpe.UserRegistration{
		FullName: name,
		Email:    fmt.Sprintf("%s-%d@provider.com", name, time.Now().UTC().UnixNano()),
		Password: userPassword,
	})
	require.NoError(t, err)
	return user
}

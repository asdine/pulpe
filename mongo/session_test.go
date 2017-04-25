package mongo_test

import (
	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/mock"
	"github.com/blankrobot/pulpe/mongo"
)

type Session struct {
	*mongo.Session
}

func (s *Session) GetAuthenticator() *mock.Authenticator {
	return s.Session.Authenticator.(*mock.Authenticator)
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

func MustGetSessions(t tester) (*Sessions, func()) {
	s := Sessions{
		NoAuth: &Session{client.Connect()},
		Red:    getSessionAs("Red"),
		Blue:   getSessionAs("Blue"),
		Green:  getSessionAs("Green"),
	}

	return &s, func() {
		// close session
		s.Close()

		_, err := client.Session.DB("").C("users").RemoveAll(nil)
		if err != nil {
			t.Error(err)
		}

		_, err = client.Session.DB("").C("boards").RemoveAll(nil)
		if err != nil {
			t.Error(err)
		}

		_, err = client.Session.DB("").C("lists").RemoveAll(nil)
		if err != nil {
			t.Error(err)
		}

		_, err = client.Session.DB("").C("cards").RemoveAll(nil)
		if err != nil {
			t.Error(err)
		}

		_, err = client.Session.DB("").C("userSessions").RemoveAll(nil)
		if err != nil {
			t.Error(err)
		}
	}
}

func getSessionAs(userID string) *Session {
	session := client.Connect()
	session.Authenticator.(*mock.Authenticator).AuthenticateFn = func(_ string) (*pulpe.User, error) {
		return &pulpe.User{ID: userID}, nil
	}

	return &Session{session}
}

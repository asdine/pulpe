package mongo_test

import (
	"os"
	"testing"

	"github.com/blankrobot/pulpe"
)

var client *Client

func MustGetSession(t tester) (pulpe.Session, func()) {
	if client == nil {
		client = MustOpenClient(t)
	}

	s := client.Connect()
	return s, func() {
		// close session
		defer s.Close()

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

func TestMain(m *testing.M) {
	code := m.Run()

	if client != nil {
		// close and remove database
		client.Close()
	}

	os.Exit(code)
}

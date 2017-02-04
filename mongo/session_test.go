package mongo_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
)

func MustGetSession(t *testing.T) (pulpe.Session, func()) {
	c := MustOpenClient(t)

	s := c.Connect()
	return s, func() {
		// close session
		defer s.Close()

		// close connection
		defer c.Close()
	}
}

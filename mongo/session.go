package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/blankrobot/pulpe"
)

var _ pulpe.Session = new(Session)

// newSession returns a new instance of Session attached to db.
func newSession(session *mgo.Session) *Session {
	s := Session{
		session: session,
		db:      session.DB(""),
	}

	s.boardService.session = &s
	return &s
}

// Session represents a connection to the database.
type Session struct {
	session *mgo.Session
	db      *mgo.Database

	now time.Time

	// Services
	cardService  CardService
	listService  ListService
	boardService BoardService
}

// CardService returns the session CardService
func (s *Session) CardService() pulpe.CardService {
	return &s.cardService
}

// ListService returns the session ListService
func (s *Session) ListService() pulpe.ListService {
	return &s.listService
}

// BoardService returns the session BoardService
func (s *Session) BoardService() pulpe.BoardService {
	return &s.boardService
}

// Close closes the mongodb session copy.
func (s *Session) Close() error {
	s.session.Close()
	return nil
}

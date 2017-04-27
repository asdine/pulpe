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
	s.boardService.store.session = &s
	s.listService.session = &s
	s.listService.store.session = &s
	s.cardService.session = &s
	s.cardService.store.session = &s
	s.userService.session = &s
	s.userSessionService.session = &s

	return &s
}

// Session represents a connection to the database.
type Session struct {
	session *mgo.Session
	db      *mgo.Database

	now time.Time

	// Services
	cardService        CardService
	listService        ListService
	boardService       BoardService
	userService        UserService
	userSessionService UserSessionService

	authenticator pulpe.Authenticator
	authToken     string
	user          *pulpe.User
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

// UserService returns the session UserService
func (s *Session) UserService() pulpe.UserService {
	return &s.userService
}

// UserSessionService returns the session UserSessionService
func (s *Session) UserSessionService() pulpe.UserSessionService {
	return &s.userSessionService
}

// Authenticate returns the current authenticate user.
func (s *Session) Authenticate() (*pulpe.User, error) {
	if s.user != nil {
		return s.user, nil
	}

	if s.authToken == "" {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	u, err := s.authenticator.Authenticate(s, s.authToken)
	if err != nil {
		return nil, err
	}

	s.user = u

	return u, nil
}

// SetAuthToken sets token as the authentication token for the session.
func (s *Session) SetAuthToken(token string) {
	s.authToken = token
}

// Close closes the mongodb session copy.
func (s *Session) Close() error {
	s.session.Close()
	return nil
}

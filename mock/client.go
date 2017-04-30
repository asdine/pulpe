package mock

import (
	"time"

	"github.com/blankrobot/pulpe"
)

// Now represents the mocked current time.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// NewClient instantiates a mock Client.
func NewClient() *Client {
	return new(Client)
}

// Client represents a mock client.
type Client struct {
	// Services
	CardService        CardService
	ListService        ListService
	BoardService       BoardService
	UserService        UserService
	UserSessionService UserSessionService
	Session            Session
}

// Connect creates a mock Session.
func (c *Client) Connect() pulpe.Session {
	c.Session.now = Now
	c.Session.boardService = &c.BoardService
	c.Session.cardService = &c.CardService
	c.Session.listService = &c.ListService
	c.Session.boardService = &c.BoardService
	c.Session.userService = &c.UserService
	c.Session.userSessionService = &c.UserSessionService
	return &c.Session
}

// Session represents a mock connection to the database.
type Session struct {
	now time.Time

	// Services
	cardService        *CardService
	listService        *ListService
	boardService       *BoardService
	userService        *UserService
	userSessionService *UserSessionService

	AuthenticateFn      func() (*pulpe.User, error)
	AuthenticateInvoked bool

	SetAuthTokenFn      func(string)
	SetAuthTokenInvoked bool
	AuthToken           string
}

// CardService returns the session CardService
func (s *Session) CardService() pulpe.CardService {
	return s.cardService
}

// ListService returns the session ListService
func (s *Session) ListService() pulpe.ListService {
	return s.listService
}

// BoardService returns the session BoardService
func (s *Session) BoardService() pulpe.BoardService {
	return s.boardService
}

// UserService returns the session UserService
func (s *Session) UserService() pulpe.UserService {
	return s.userService
}

// UserSessionService returns the session UserSessionService
func (s *Session) UserSessionService() pulpe.UserSessionService {
	return s.userSessionService
}

// Authenticate runs AuthenticateFn and sets AuthenticateInvoked to true when invoked.
func (s *Session) Authenticate() (*pulpe.User, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn()
}

// SetAuthToken sets SetAuthTokenInvoked to true when invoked and runs SetAuthTokenFn if it exists,
// otherwhise it stores the token in AuthToken.
func (s *Session) SetAuthToken(token string) {
	s.SetAuthTokenInvoked = true
	if s.SetAuthTokenFn != nil {
		s.SetAuthTokenFn(token)
	} else {
		s.AuthToken = token
	}
}

// Close session.
func (s *Session) Close() error {
	return nil
}

// Authenticator represents a mock Authenticator.
type Authenticator struct {
	AuthenticateFn      func(pulpe.Session, string) (*pulpe.User, error)
	AuthenticateInvoked bool
}

// Authenticate runs AuthenticateFn and sets AuthenticateInvoked to true when invoked.
func (a *Authenticator) Authenticate(session pulpe.Session, token string) (*pulpe.User, error) {
	a.AuthenticateInvoked = true
	if a.AuthenticateFn == nil {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	return a.AuthenticateFn(session, token)
}

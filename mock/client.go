package mock

import (
	"time"

	"github.com/blankrobot/pulpe"
)

// Now represents the mocked current time.
var Now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// NewClient instantiates a mock Client.
func NewClient() *Client {
	return &Client{}
}

// Client represents a mock client.
type Client struct {
	// Services
	CardService  CardService
	ListService  ListService
	BoardService BoardService
	UserService  UserService
}

// Connect creates mock Session.
func (c *Client) Connect() pulpe.Session {
	return &Session{
		now:          Now,
		cardService:  &c.CardService,
		listService:  &c.ListService,
		boardService: &c.BoardService,
		userService:  &c.UserService,
	}
}

// Session represents a mock connection to the database.
type Session struct {
	now time.Time

	// Services
	cardService  *CardService
	listService  *ListService
	boardService *BoardService
	userService  *UserService
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

// Close session.
func (s *Session) Close() error {
	return nil
}

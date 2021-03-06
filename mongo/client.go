package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/blankrobot/pulpe"
)

// NewClient instantiates a new Client.
func NewClient(uri string) *Client {
	return &Client{
		Now: time.Now,
		URI: uri,
	}
}

// Client represents a client to the underlying MongoDB database.
type Client struct {
	// MongoDB database uri.
	URI string

	// Returns the current time.
	Now func() time.Time

	// Authenticator
	Authenticator pulpe.Authenticator

	Session *mgo.Session
}

// Open opens and initializes the MongoDB database.
func (c *Client) Open() error {
	var err error

	c.Session, err = mgo.Dial(c.URI)
	if err != nil {
		return err
	}

	return c.EnsureIndexes()
}

// EnsureIndexes creates indexes if they don't exist.
func (c *Client) EnsureIndexes() error {
	session := c.Connect()
	defer session.Close()

	err := session.CardService().(*CardService).ensureIndexes()
	if err != nil {
		return err
	}

	err = session.ListService().(*ListService).ensureIndexes()
	if err != nil {
		return err
	}

	err = session.BoardService().(*BoardService).ensureIndexes()
	if err != nil {
		return err
	}

	err = session.UserService().(*UserService).ensureIndexes()
	if err != nil {
		return err
	}

	return session.UserSessionService().(*UserSessionService).ensureIndexes()
}

// Close closes then underlying MongoDB database.
func (c *Client) Close() error {
	c.Session.Close()

	return nil
}

// Connect creates a new session.
func (c *Client) Connect() pulpe.Session {
	s := newSession(c.Session.Copy())
	s.now = c.Now().UTC()
	s.authenticator = c.Authenticator
	return s
}

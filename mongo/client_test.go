package mongo_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/blankrobot/pulpe/mock"
	"github.com/blankrobot/pulpe/mongo"
)

var client *Client

func TestMain(m *testing.M) {
	var err error

	client, err = OpenClient()
	if err != nil {
		log.Fatal(err)
	}
	defer sessions.Close()
	defer client.Close()

	code := m.Run()

	os.Exit(code)
}

// Client is a test wrapper for mongo.Client.
type Client struct {
	*mongo.Client
}

// NewClient returns a new instance of Client.
func NewClient(uri string) *Client {
	c := Client{
		Client: mongo.NewClient(fmt.Sprintf("%s/pulpe-tests", uri)),
	}
	c.Client.Now = Now
	c.Client.Authenticator = new(mock.Authenticator)

	return &c
}

// OpenClient returns an new, open instance of Client.
func OpenClient() (*Client, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	c := NewClient(uri)
	if err := c.Open(); err != nil {
		return nil, err
	}

	return c, nil
}

// Connect creates a new session.
func (c *Client) Connect() *mongo.Session {
	c.Client.Authenticator = new(mongo.Authenticator)
	return c.Client.Connect().(*mongo.Session)
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	err := c.Session.DB("").DropDatabase()
	if err != nil {
		return err
	}

	return c.Client.Close()
}

func Now() time.Time {
	n := time.Now().UTC()
	return time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), 0, time.UTC)
}

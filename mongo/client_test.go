package mongo_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/blankrobot/pulpe/mock"
	"github.com/blankrobot/pulpe/mongo"
)

// Client is a test wrapper for mongo.Client.
type Client struct {
	*mongo.Client
}

// NewClient returns a new instance of Client.
func NewClient(uri string) *Client {
	c := Client{
		Client: mongo.NewClient(fmt.Sprintf("%s/pulpe-tests-%d", uri, time.Now().UnixNano())),
	}
	c.Now = func() time.Time { return mock.Now }

	return &c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient(t *testing.T) *Client {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		t.Skip("Missing MONGO_URI environment variable.")
	}

	c := NewClient(uri)
	if err := c.Open(); err != nil {
		t.Error(err)
	}

	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	err := c.Session.DB("").DropDatabase()
	if err != nil {
		return err
	}

	return c.Client.Close()
}

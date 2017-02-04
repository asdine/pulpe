package mongo_test

import (
	"encoding/json"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

var settings = json.RawMessage([]byte("{}"))

// Ensure boards can be created and retrieved.
func TestBoardService_CreateBoard(t *testing.T) {
	c := MustOpenClient(t)
	defer c.Close()

	s := c.Connect().BoardService()

	t.Run("New", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "XXX",
			Settings: &settings,
		}
		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)

		// Retrieve board and compare.
		other, err := s.Board(board.ID)
		require.NoError(t, err)
		require.Equal(t, board, other)
	})
}

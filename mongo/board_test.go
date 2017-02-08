package mongo_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

var settings = json.RawMessage([]byte("{}"))

// Ensure boards can be created and retrieved.
func TestBoardService_CreateBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

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

func TestBoardService_Board(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("Exists", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "ZZZ",
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

	t.Run("Not found", func(t *testing.T) {
		// Trying to fetch a board that doesn't exist.
		_, err := s.Board("QQQ")
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})
}

func TestBoardService_Boards(t *testing.T) {
	t.Parallel()

	t.Run("Exists", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()

		s := session.BoardService()

		for i := 0; i < 5; i++ {
			b := pulpe.BoardCreate{
				Name:     fmt.Sprintf("board%d", i),
				Settings: &settings,
			}
			// Create new board.
			_, err := s.CreateBoard(&b)
			require.NoError(t, err)
		}

		// Retrieve boards.
		boards, err := s.Boards()
		require.NoError(t, err)
		require.Len(t, boards, 5)
		require.Equal(t, boards[0].Name, "board0")
		require.Equal(t, boards[4].Name, "board4")
	})

	t.Run("No boards", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()

		s := session.BoardService()

		boards, err := s.Boards()
		require.NoError(t, err)
		require.Len(t, boards, 0)
	})
}

func TestBoardService_DeleteBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("Exists", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "Board1",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)

		// Delete board.
		err = s.DeleteBoard(board.ID)
		require.NoError(t, err)

		_, err = s.Board(board.ID)
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to delete a board that doesn't exist.
		err := s.DeleteBoard("QQQ")
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})
}

func TestBoardService_UpdateBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("Exists", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "name",
			Settings: &settings,
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)

		// Update a single field.
		newName := "new name"
		newSettings := json.RawMessage([]byte(`{"a":"b"}`))
		updatedBoard, err := s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name:     &newName,
			Settings: &newSettings,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedBoard)

		// Retrieve board and check.
		other, err := s.Board(board.ID)
		require.NoError(t, err)
		require.Equal(t, newName, other.Name)
		require.NotNil(t, other.UpdatedAt)
		require.Equal(t, other, updatedBoard)
		// Set zero values.
		newName = ""
		newSettings = []byte("")
		updatedBoard, err = s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name:     &newName,
			Settings: &newSettings,
		})
		require.NoError(t, err)
		require.Zero(t, updatedBoard.Name)
		require.Nil(t, updatedBoard.Settings)

		// Retrieve board and check.
		other, err = s.Board(board.ID)
		require.NoError(t, err)
		require.Zero(t, other.Name)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a board that doesn't exist with no patch.
		updatedBoard, err := s.UpdateBoard("QQQ", &pulpe.BoardUpdate{})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)

		// Trying to update a board that doesn't exist with a patch.
		newName := "new name"
		updatedBoard, err = s.UpdateBoard("QQQ", &pulpe.BoardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)
	})
}

package mongo_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

var settings = json.RawMessage([]byte("{}"))

func newBoardID() string {
	return bson.NewObjectId().Hex()
}

// Ensure boards can be created and retrieved.
func TestBoardService_CreateBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("New", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "XXX YYY ",
			Settings: &settings,
		}
		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "xxx-yyy")

		// Retrieve board and compare.
		other, err := s.Board(board.ID)
		require.NoError(t, err)
		require.Equal(t, board, other)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "ZZZ KK ",
			Settings: &settings,
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "zzz-kk")

		// Create second board with slightly different name that generates the same slug.
		b = pulpe.BoardCreate{
			Name:     "  ZZZ   KK ",
			Settings: &settings,
		}
		_, err = s.CreateBoard(&b)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardExists, err)
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
		id := bson.NewObjectId().Hex()
		_, err := s.Board(id)
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

		// Retrieve boards without filter.
		boards, err := s.Boards(nil)
		require.NoError(t, err)
		require.Len(t, boards, 5)
		require.Equal(t, boards[0].Name, "board0")
		require.Equal(t, boards[4].Name, "board4")

		// Retrieve boards with filter.
		boards, err = s.Boards(map[string]string{
			"slug": "board4",
		})
		require.NoError(t, err)
		require.Len(t, boards, 1)
		require.Equal(t, boards[0].Name, "board4")
	})

	t.Run("No boards", func(t *testing.T) {
		session, cleanup := MustGetSession(t)
		defer cleanup()

		s := session.BoardService()

		boards, err := s.Boards(nil)
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
		id := bson.NewObjectId().Hex()
		err := s.DeleteBoard(id)
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
		require.Equal(t, "new-name", other.Slug)
		require.Equal(t, other, updatedBoard)

		// Set zero values.
		newSettings = []byte("")
		updatedBoard, err = s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Settings: &newSettings,
		})
		require.NoError(t, err)
		require.Nil(t, updatedBoard.Settings)

		// Update a single field with the same value.
		updatedBoard, err = s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedBoard)

		// Retrieve board and check.
		other, err = s.Board(board.ID)
		require.NoError(t, err)
		require.Zero(t, other.Settings)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a board that doesn't exist with no patch.
		id := bson.NewObjectId().Hex()
		updatedBoard, err := s.UpdateBoard(id, &pulpe.BoardUpdate{})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)

		// Trying to update a board that doesn't exist with a patch.
		newName := "new name 2"
		updatedBoard, err = s.UpdateBoard(id, &pulpe.BoardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)
	})
}

func BenchmarkCreateBoard(b *testing.B) {
	session, cleanup := MustGetSession(b)
	defer cleanup()

	s := session.BoardService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc := pulpe.BoardCreate{
			Name: fmt.Sprintf("name%d", i),
		}
		_, err := s.CreateBoard(&bc)
		require.NoError(b, err)
	}
}

func BenchmarkGetBoard(b *testing.B) {
	session, cleanup := MustGetSession(b)
	defer cleanup()

	s := session.BoardService()

	var id string
	for i := 0; i < 1000; i++ {
		bc := pulpe.BoardCreate{
			Name: fmt.Sprintf("name%d", i),
		}
		board, err := s.CreateBoard(&bc)
		require.NoError(b, err)
		id = board.ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := s.Board(id)
		require.NoError(b, err)
	}
}

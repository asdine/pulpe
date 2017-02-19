package mongo_test

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

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
			Name: "XXX YYY ",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "xxx-yyy")

		// Retrieve board and compare.
		other, err := s.Board("xxx-yyy")
		require.NoError(t, err)
		require.Equal(t, board, other)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "ZZZ KK ",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "zzz-kk")

		// Create second board with slightly different name that generates the same slug.
		b = pulpe.BoardCreate{
			Name: "  ZZZ   KK ",
		}
		board, err = s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", board.Slug)
	})
}

func TestBoardService_Board(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("Exists", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "ZZZ",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)

		// Retrieve board and compare.
		other, err := s.Board(board.Slug)
		require.NoError(t, err)
		require.Equal(t, board, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to fetch a board that doesn't exist.
		_, err := s.Board("something")
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
				Name: fmt.Sprintf("board%d", i),
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

		_, err = s.Board(board.Slug)
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to delete a board that doesn't exist.
		err := s.DeleteBoard("something")
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})
}

func TestBoardService_UpdateBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.BoardService()

	t.Run("OK", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "name",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)

		// Update a single field.
		newName := "new name"
		updatedBoard, err := s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name: &newName,
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

		// Update a single field with the same value.
		updatedBoard, err = s.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.Equal(t, "new-name", updatedBoard.Slug)
		require.NotNil(t, updatedBoard)

		// Retrieve board and check.
		other, err = s.Board(board.ID)
		require.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a board that doesn't exist with no patch.
		updatedBoard, err := s.UpdateBoard(newBoardID(), &pulpe.BoardUpdate{})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)

		// Trying to update a board that doesn't exist with a patch.
		newName := "new name 2"
		updatedBoard, err = s.UpdateBoard(newBoardID(), &pulpe.BoardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		b1 := pulpe.BoardCreate{
			Name: "hello",
		}

		b2 := pulpe.BoardCreate{
			Name: "goodbye",
		}

		// Create board 1.
		_, err := s.CreateBoard(&b1)
		require.NoError(t, err)

		// Create board 2.
		board2, err := s.CreateBoard(&b2)
		require.NoError(t, err)

		// Update b2.
		newName := "hello"
		updatedBoard, err := s.UpdateBoard(board2.ID, &pulpe.BoardUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.Equal(t, "hello-1", updatedBoard.Slug)
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

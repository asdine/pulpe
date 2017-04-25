package mongo_test

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func newBoardID() string {
	return bson.NewObjectId().Hex()
}

func newBoard(t require.TestingT, session *Session) *pulpe.Board {
	board, err := session.BoardService().CreateBoard(&pulpe.BoardCreation{
		Name: fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
	})
	require.NoError(t, err)

	return board
}

// Ensure boards can be created and retrieved.
func TestBoardService_CreateBoard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.BoardService()
		b := pulpe.BoardCreation{
			Name: "XXX YYY ",
		}

		// Create new board.
		_, err := s.CreateBoard(&b)
		require.Error(t, err)
		require.True(t, sessions.NoAuth.GetAuthenticator().AuthenticateInvoked)
	})

	t.Run("New", func(t *testing.T) {
		s := sessions.Red.BoardService()

		b := pulpe.BoardCreation{
			Name: "XXX YYY ",
		}

		// Create new board.
		board, err := s.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, "xxx-yyy", board.Slug)

		// Retrieve board and compare.
		other, err := s.Board(board.ID)
		require.NoError(t, err)
		require.Equal(t, board, other)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		s1 := sessions.Red.BoardService()
		s2 := sessions.Blue.BoardService()

		b := pulpe.BoardCreation{
			Name: "ZZZ KK ",
		}

		// Create new board.
		board, err := s1.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "zzz-kk")

		// Create another board with the same slug but with another user.
		board, err = s2.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, board.Slug, "zzz-kk")

		// Create second board with slightly different name that generates the same slug.
		b = pulpe.BoardCreation{
			Name: "  ZZZ   KK ",
		}

		board, err = s1.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", board.Slug)

		board, err = s2.CreateBoard(&b)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", board.Slug)
	})
}

func TestBoardService_Board(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.BoardService()

		// Get a board.
		_, err := s.Board("someid")
		require.Error(t, err)
		require.True(t, sessions.NoAuth.GetAuthenticator().AuthenticateInvoked)
	})

	t.Run("Exists", func(t *testing.T) {
		s1 := sessions.Red.BoardService()
		s2 := sessions.Blue.BoardService()

		b := pulpe.BoardCreation{
			Name: "ZZZ",
		}

		// Create new board as red
		board1, err := s1.CreateBoard(&b)
		require.NoError(t, err)

		// Create new board as blue
		board2, err := s2.CreateBoard(&b)
		require.NoError(t, err)

		// Retrieve board1 as blue
		other, err := s2.Board(board1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardNotFound, err)

		// Retrieve board2 as blue
		other, err = s2.Board(board2.ID)
		require.NoError(t, err)
		require.Equal(t, board2, other)

		// Retrieve board2 as red
		other, err = s1.Board(board2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardNotFound, err)

		// Retrieve board1 as red
		other, err = s1.Board(board1.ID)
		require.NoError(t, err)
		require.Equal(t, board1, other)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Red.BoardService()
		// Trying to fetch a board that doesn't exist.
		_, err := s.Board("something")
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})
}

func TestBoardService_Boards(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.BoardService()

		_, err := s.Boards()
		require.Error(t, err)
		require.True(t, sessions.NoAuth.GetAuthenticator().AuthenticateInvoked)
	})

	t.Run("Exists", func(t *testing.T) {
		s1 := sessions.Red.BoardService()
		s2 := sessions.Blue.BoardService()

		for i := 0; i < 5; i++ {
			b := pulpe.BoardCreation{
				Name: fmt.Sprintf("board%d", i),
			}
			// Create new board as red.
			_, err := s1.CreateBoard(&b)
			require.NoError(t, err)

			_, err = s2.CreateBoard(&b)
			require.NoError(t, err)
		}

		// Retrieve boards.
		boards1, err := s1.Boards()
		require.NoError(t, err)
		require.Len(t, boards1, 5)
		require.Equal(t, boards1[0].Name, "board0")
		require.Equal(t, boards1[4].Name, "board4")

		// Retrieve boards.
		boards2, err := s2.Boards()
		require.NoError(t, err)
		require.Len(t, boards2, 5)
		require.Equal(t, boards2[0].Name, "board0")
		require.Equal(t, boards2[4].Name, "board4")

		require.NotEmpty(t, boards1, boards2)
	})

	t.Run("No boards", func(t *testing.T) {
		s := sessions.Green.BoardService()

		boards, err := s.Boards()
		require.NoError(t, err)
		require.Len(t, boards, 0)
	})
}

func TestBoardService_DeleteBoard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.BoardService()

		// Trying to delete a board.
		err := s.DeleteBoard("something")
		require.Error(t, err)
		require.True(t, sessions.NoAuth.GetAuthenticator().AuthenticateInvoked)
	})

	t.Run("Exists", func(t *testing.T) {
		s1 := sessions.Red.BoardService()
		s2 := sessions.Blue.BoardService()

		b := pulpe.BoardCreation{
			Name: "Board",
		}

		// Create new board as red
		board1, err := s1.CreateBoard(&b)
		require.NoError(t, err)

		// Create new board as blue
		board2, err := s2.CreateBoard(&b)
		require.NoError(t, err)

		// Delete board1 as blue
		err = s2.DeleteBoard(board1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardNotFound, err)

		// Delete board2 as red
		err = s1.DeleteBoard(board2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardNotFound, err)

		// Delete board1 as red
		err = s1.DeleteBoard(board1.ID)
		require.NoError(t, err)

		// Delete board2 as blue
		err = s2.DeleteBoard(board2.ID)
		require.NoError(t, err)

		_, err = s1.Board(board1.Slug)
		require.Equal(t, pulpe.ErrBoardNotFound, err)

		_, err = s2.Board(board2.Slug)
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.BoardService()

		// Trying to delete a board that doesn't exist.
		err := s.DeleteBoard("something")
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})
}

func TestBoardService_UpdateBoard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.BoardService()

		// Trying to delete a board.
		err := s.DeleteBoard("something")
		require.Error(t, err)
		require.True(t, sessions.NoAuth.GetAuthenticator().AuthenticateInvoked)
	})

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.BoardService()

		b := pulpe.BoardCreation{
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

	t.Run("Bad user", func(t *testing.T) {
		s1 := sessions.Red.BoardService()
		s2 := sessions.Blue.BoardService()

		b := pulpe.BoardCreation{
			Name: "name",
		}

		// Create new board as red.
		board, err := s1.CreateBoard(&b)
		require.NoError(t, err)

		// Update as blue.
		newName := "new name"
		_, err = s2.UpdateBoard(board.ID, &pulpe.BoardUpdate{
			Name: &newName,
		})
		require.Error(t, err)
		require.Equal(t, pulpe.ErrBoardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Red.BoardService()
		// Trying to update a board that doesn't exist with no patch.
		updatedBoard, err := s.UpdateBoard(newBoardID(), new(pulpe.BoardUpdate))
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)

		// Trying to update a board that doesn't exist with a patch.
		newName := "new name 2"
		updatedBoard, err = s.UpdateBoard(newBoardID(), &pulpe.BoardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrBoardNotFound, err)
		require.Nil(t, updatedBoard)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		s := sessions.Green.BoardService()

		b1 := pulpe.BoardCreation{
			Name: "hello",
		}

		b2 := pulpe.BoardCreation{
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
	sessions, cleanup := MustGetSessions(b)
	defer cleanup()

	s := sessions.Red.BoardService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc := pulpe.BoardCreation{
			Name: fmt.Sprintf("name%d", i),
		}
		_, err := s.CreateBoard(&bc)
		require.NoError(b, err)
	}
}

func BenchmarkGetBoard(b *testing.B) {
	sessions, cleanup := MustGetSessions(b)
	defer cleanup()

	s := sessions.Red.BoardService()

	var id string
	for i := 0; i < 1000; i++ {
		bc := pulpe.BoardCreation{
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

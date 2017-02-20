package mongo_test

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func newListID() string {
	return bson.NewObjectId().Hex()
}

func TestListService_CreateList(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("New", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: newBoardID(),
			Name:    "Name",
		}

		// Create new list.
		list, err := s.CreateList(&l)
		require.NoError(t, err)
		require.NotZero(t, list.ID)
		require.Equal(t, list.Name, l.Name)

		// Retrieve list and compare.
		other, err := s.List(list.ID)
		require.NoError(t, err)
		require.Equal(t, list.ID, other.ID)
		require.Equal(t, list.Name, other.Name)
		require.Equal(t, list.BoardID, other.BoardID)
		require.Equal(t, list.CreatedAt.UTC(), other.CreatedAt.UTC())
		require.Equal(t, list.UpdatedAt, other.UpdatedAt)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		boardID := newBoardID()
		l := pulpe.ListCreate{
			Name:    "ZZZ KK ",
			BoardID: boardID,
		}

		// Create new list.
		list, err := s.CreateList(&l)
		require.NoError(t, err)
		require.Equal(t, list.Slug, "zzz-kk")

		// Create second list with slightly different name that generates the same slug.
		l = pulpe.ListCreate{
			Name:    "  ZZZ   KK ",
			BoardID: boardID,
		}
		list, err = s.CreateList(&l)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", list.Slug)
	})
}

func TestListService_List(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: newBoardID(),
		}

		// Create new list.
		list, err := s.CreateList(&l)
		require.NoError(t, err)

		// Retrieve list and compare.
		other, err := s.List(list.ID)
		require.NoError(t, err)
		require.Equal(t, list.ID, other.ID)
		require.Equal(t, list.Name, other.Name)
		require.Equal(t, list.BoardID, other.BoardID)
		require.Equal(t, list.CreatedAt.UTC(), other.CreatedAt.UTC())
		require.Equal(t, list.UpdatedAt, other.UpdatedAt)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to fetch a list that doesn't exist.
		_, err := s.List(newListID())
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

func TestListService_DeleteList(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: newBoardID(),
		}

		// Create new list.
		list, err := s.CreateList(&l)
		require.NoError(t, err)

		// Delete list.
		err = s.DeleteList(list.ID)
		require.NoError(t, err)

		_, err = s.List(list.ID)
		require.Equal(t, pulpe.ErrListNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to delete a list that doesn't exist.
		err := s.DeleteList(newListID())
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

func TestListService_DeleteListsByBoardID(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		board1 := newBoardID()
		board2 := newBoardID()

		for i := 0; i < 10; i++ {
			var c pulpe.ListCreate

			if i%2 != 0 {
				c.BoardID = board1
			} else {
				c.BoardID = board2
			}

			// Create new list.
			_, err := s.CreateList(&c)
			require.NoError(t, err)
		}

		// Delete by board id.
		err := s.DeleteListsByBoardID(board1)
		require.NoError(t, err)

		lists, err := s.ListsByBoard(board1)
		require.NoError(t, err)
		require.Len(t, lists, 0)
	})

	t.Run("Not found", func(t *testing.T) {
		// Calling with a boardID with no associated lists.
		err := s.DeleteListsByBoardID(newBoardID())
		require.NoError(t, err)
	})
}

func TestListService_UpdateList(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		c := pulpe.ListCreate{
			BoardID: newBoardID(),
			Name:    "name",
		}

		// Create new list.
		list, err := s.CreateList(&c)
		require.NoError(t, err)

		// Update a single field.
		newName := "new name"
		updatedList, err := s.UpdateList(list.ID, &pulpe.ListUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedList)

		// Retrieve list and check.
		other, err := s.List(list.ID)
		require.NoError(t, err)
		require.Equal(t, newName, other.Name)
		require.NotNil(t, other.UpdatedAt)
		require.Equal(t, updatedList, other)

		// Set zero values.
		newName = ""
		updatedList, err = s.UpdateList(list.ID, &pulpe.ListUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedList)

		// Retrieve list and check.
		other, err = s.List(list.ID)
		require.NoError(t, err)
		require.Zero(t, other.Name)
		require.Equal(t, updatedList, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a list that doesn't exist with no patch.
		updatedList, err := s.UpdateList(newListID(), &pulpe.ListUpdate{})
		require.Equal(t, pulpe.ErrListNotFound, err)
		require.Nil(t, updatedList)

		// Trying to update a list that doesn't exist with a patch.
		newName := "new name"
		updatedList, err = s.UpdateList(newListID(), &pulpe.ListUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrListNotFound, err)
		require.Nil(t, updatedList)
	})
}

func TestListService_ListsByBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		boardID1 := newBoardID()
		boardID2 := newBoardID()
		for i := 0; i < 6; i++ {
			var l pulpe.ListCreate

			if i%2 == 0 {
				l.BoardID = boardID1
			} else {
				l.BoardID = boardID2
			}

			// Create new list.
			_, err := s.CreateList(&l)
			require.NoError(t, err)
		}

		lists, err := s.ListsByBoard(boardID1)
		require.NoError(t, err)
		require.Len(t, lists, 3)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to find lists of a board that doesn't exist.
		lists, err := s.ListsByBoard(newBoardID())
		require.NoError(t, err)
		require.Len(t, lists, 0)
	})
}

package mongo_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func TestListService_CreateList(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("New", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: "BoardID",
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

	t.Run("No BoardID", func(t *testing.T) {
		// Trying to create a card with no ID.
		var l pulpe.ListCreate

		_, err := s.CreateList(&l)
		require.Equal(t, pulpe.ErrListBoardIDRequired, err)
	})
}

func TestListService_List(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: "BoardID",
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
		_, err := s.List("QQQ")
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

func TestListService_DeleteList(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.ListService()

	t.Run("Exists", func(t *testing.T) {
		l := pulpe.ListCreate{
			BoardID: "BoardID",
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
		err := s.DeleteList("QQQ")
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

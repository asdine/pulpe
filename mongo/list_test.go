package mongo_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

func newListID() string {
	return bson.NewObjectId().Hex()
}

func newList(t require.TestingT, session *Session) *pulpe.List {
	return newListWithBoardID(t, session, newBoard(t, session).ID)
}

func newListWithBoardID(t require.TestingT, session *Session, boardID string) *pulpe.List {
	list, err := session.ListService().CreateList(boardID, &pulpe.ListCreation{
		Name: fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
	})
	require.NoError(t, err)

	return list
}

func TestListService_CreateList(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.ListService()
		l := pulpe.ListCreation{
			Name:     "Name",
			Position: 1000,
		}

		// Create new list.
		_, err := s.CreateList(newBoardID(), &l)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("New", func(t *testing.T) {
		s := sessions.Red.ListService()
		l := pulpe.ListCreation{
			Name:     "Name",
			Position: 1000,
		}

		// Create new list.
		list, err := s.CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)
		require.NotZero(t, list.ID)
		require.Equal(t, list.Name, l.Name)
		require.Equal(t, list.Position, l.Position)
		require.Equal(t, "name", list.Slug)

		// Retrieve list and compare.
		other, err := s.List(list.ID)
		require.NoError(t, err)
		require.Equal(t, list.ID, other.ID)
		require.Equal(t, list.Name, other.Name)
		require.Equal(t, list.Position, other.Position)
		require.Equal(t, list.BoardID, other.BoardID)
		require.Equal(t, list.CreatedAt.UTC(), other.CreatedAt.UTC())
		require.Equal(t, list.UpdatedAt, other.UpdatedAt)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		s := sessions.Red.ListService()

		boardID := newBoard(t, sessions.Red).ID
		l := pulpe.ListCreation{
			Name: "ZZZ KK ",
		}

		// Create new list.
		list, err := s.CreateList(boardID, &l)
		require.NoError(t, err)
		require.Equal(t, list.Slug, "zzz-kk")

		// Create second list with slightly different name that generates the same slug.
		l = pulpe.ListCreation{
			Name: "  ZZZ   KK ",
		}
		list, err = s.CreateList(boardID, &l)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", list.Slug)
	})
}

func TestListService_List(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.ListService()

		// Get a list
		_, err := s.List("someid")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s1 := sessions.Red.ListService()
		s2 := sessions.Blue.ListService()

		l := pulpe.ListCreation{
			Name: "ZZZ",
		}

		// Create new list as red
		list1, err := s1.CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Create new list as blue
		list2, err := s2.CreateList(newBoard(t, sessions.Blue).ID, &l)
		require.NoError(t, err)

		// Retrieve list1 as blue
		other, err := s2.List(list1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)

		// Retrieve list2 as blue
		other, err = s2.List(list2.ID)
		require.NoError(t, err)
		require.Equal(t, list2, other)

		// Retrieve list2 as red
		other, err = s1.List(list2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)

		// Retrieve list1 as red
		other, err = s1.List(list1.ID)
		require.NoError(t, err)
		require.Equal(t, list1, other)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.ListService()

		// Trying to fetch a list that doesn't exist.
		_, err := s.List(newListID())
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

func TestListService_DeleteList(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.ListService()

		// Trying to delete a list.
		err := s.DeleteList("something")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s1 := sessions.Red.ListService()
		s2 := sessions.Blue.ListService()

		l := pulpe.ListCreation{
			Name: "name",
		}

		// Create new list as red
		list1, err := s1.CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Create new list as blue
		list2, err := s2.CreateList(newBoard(t, sessions.Blue).ID, &l)
		require.NoError(t, err)

		// Delete list1 as blue
		err = s2.DeleteList(list1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)

		// Delete list2 as red
		err = s1.DeleteList(list2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)

		// Delete list1 as red
		err = s1.DeleteList(list1.ID)
		require.NoError(t, err)

		// Delete list2 as blue
		err = s2.DeleteList(list2.ID)
		require.NoError(t, err)

		_, err = s1.List(list1.Slug)
		require.Equal(t, pulpe.ErrListNotFound, err)

		_, err = s2.List(list2.Slug)
		require.Equal(t, pulpe.ErrListNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.ListService()

		// Trying to delete a list that doesn't exist.
		err := s.DeleteList(newListID())
		require.Equal(t, pulpe.ErrListNotFound, err)
	})
}

func TestListService_DeleteListsByBoardID(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.ListService()

		board1 := newBoard(t, sessions.Red).ID
		board2 := newBoard(t, sessions.Red).ID

		var err error
		for i := 0; i < 10; i++ {
			var c pulpe.ListCreation

			if i%2 != 0 {
				_, err = s.CreateList(board1, &c)
			} else {
				_, err = s.CreateList(board2, &c)
			}

			require.NoError(t, err)
		}

		// Delete by board id.
		err = s.DeleteListsByBoardID(board1)
		require.NoError(t, err)

		lists, err := s.ListsByBoard(board1)
		require.NoError(t, err)
		require.Len(t, lists, 0)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.ListService()

		// Calling with a boardID with no associated lists.
		err := s.DeleteListsByBoardID(newBoard(t, sessions.Green).ID)
		require.NoError(t, err)
	})
}

func TestListService_UpdateList(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.ListService()

		// Update a list.
		_, err := s.UpdateList("someid", new(pulpe.ListUpdate))
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.ListService()

		l := pulpe.ListCreation{
			Name:     "name",
			Position: 1000,
		}

		// Create new list.
		list, err := s.CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Update a single field.
		newName := "new name"
		updatedList, err := s.UpdateList(list.ID, &pulpe.ListUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedList)

		// Update multiple fields.
		newName = "other name"
		newPos := 500.0
		updatedList, err = s.UpdateList(list.ID, &pulpe.ListUpdate{
			Name:     &newName,
			Position: &newPos,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedList)

		// Retrieve list and check.
		other, err := s.List(list.ID)
		require.NoError(t, err)
		require.Equal(t, newName, other.Name)
		require.Equal(t, newPos, other.Position)
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

	t.Run("Bad user", func(t *testing.T) {
		s1 := sessions.Red.ListService()
		s2 := sessions.Blue.ListService()

		l := pulpe.ListCreation{
			Name: "name",
		}

		// Create new List as red.
		list, err := s1.CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Update as blue.
		newName := "new name"
		_, err = s2.UpdateList(list.ID, &pulpe.ListUpdate{
			Name: &newName,
		})
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.ListService()

		// Trying to update a list that doesn't exist with no patch.
		updatedList, err := s.UpdateList(newListID(), new(pulpe.ListUpdate))
		require.Equal(t, pulpe.ErrListNotFound, err)
		require.Nil(t, updatedList)

		// Trying to update a list that doesn't exist with a patch.
		newName := "new name"
		updatedList, err = s.UpdateList(newListID(), &pulpe.ListUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrListNotFound, err)
		require.Nil(t, updatedList)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		s := sessions.Green.ListService()

		boardID := newBoard(t, sessions.Green).ID
		l1 := pulpe.ListCreation{
			Name: "hello",
		}

		l2 := pulpe.ListCreation{
			Name: "goodbye",
		}

		// Create List 1.
		_, err := s.CreateList(boardID, &l1)
		require.NoError(t, err)

		// Create List 2.
		list2, err := s.CreateList(boardID, &l2)
		require.NoError(t, err)

		// Update l2.
		newName := "hello"
		updatedList, err := s.UpdateList(list2.ID, &pulpe.ListUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.Equal(t, "hello-1", updatedList.Slug)
	})
}

func TestListService_ListsByBoard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.ListService()

		boardID1 := newBoard(t, sessions.Red).ID
		boardID2 := newBoard(t, sessions.Red).ID

		var err error
		for i := 0; i < 6; i++ {
			var l pulpe.ListCreation

			if i%2 == 0 {
				_, err = s.CreateList(boardID1, &l)
			} else {
				_, err = s.CreateList(boardID2, &l)
			}

			// Create new list.

			require.NoError(t, err)
		}

		lists, err := s.ListsByBoard(boardID1)
		require.NoError(t, err)
		require.Len(t, lists, 3)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.ListService()

		// Trying to find lists of a board that doesn't exist.
		lists, err := s.ListsByBoard(newBoard(t, sessions.Green).ID)
		require.NoError(t, err)
		require.Len(t, lists, 0)
	})
}

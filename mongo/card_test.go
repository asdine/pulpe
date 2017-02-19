package mongo_test

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func newCardID() string {
	return bson.NewObjectId().Hex()
}

// Ensure cards can be created and retrieved.
func TestCardService_CreateCard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("New", func(t *testing.T) {
		c := pulpe.CardCreate{
			Name:        "YYY",
			ListID:      newListID(),
			BoardID:     newBoardID(),
			Description: "MY CARD",
			Position:    1,
		}
		// Create new card.
		card, err := s.CreateCard(&c)
		require.NoError(t, err)
		require.NotZero(t, card.ID)
		require.Equal(t, "yyy", card.Slug)

		// Retrieve card and compare.
		other, err := s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, card, other)
	})
}

// Ensure cards can be retrieved.
func TestCardService_Card(t *testing.T) {
	t.Parallel()
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("OK", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:  newListID(),
			BoardID: newBoardID(),
			Name:    "name",
		}

		// Create new card.
		card, err := s.CreateCard(&c)
		require.NoError(t, err)

		// Retrieve card and compare.
		other, err := s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, card, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to fetch a card that doesn't exist.
		_, err := s.Card("something")
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCard(t *testing.T) {
	t.Parallel()
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("OK", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:  newListID(),
			BoardID: newBoardID(),
			Name:    "name",
		}

		// Create new card.
		card, err := s.CreateCard(&c)
		require.NoError(t, err)

		// Delete card.
		err = s.DeleteCard(card.ID)
		require.NoError(t, err)

		// Try to delete the same card.
		err = s.DeleteCard(card.ID)
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to delete a card that doesn't exist.
		err := s.DeleteCard(bson.NewObjectId().Hex())
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCardsByListID(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("OK", func(t *testing.T) {
		list1 := newListID()
		list2 := newListID()
		boardID := newBoardID()

		for i := 0; i < 10; i++ {
			c := pulpe.CardCreate{
				BoardID: boardID,
				Name:    "name",
			}

			if i%2 != 0 {
				c.ListID = list1
			} else {
				c.ListID = list2
			}

			// Create new card.
			_, err := s.CreateCard(&c)
			require.NoError(t, err)
		}

		// Delete card.
		err := s.DeleteCardsByListID(list1)
		require.NoError(t, err)

		cards, err := s.CardsByBoard(boardID)
		require.NoError(t, err)
		require.Len(t, cards, 5)
		for _, card := range cards {
			require.Equal(t, list2, card.ListID)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		// Calling with a listID with no associated cards.
		err := s.DeleteCardsByListID(newListID())
		require.NoError(t, err)
	})
}

func TestCardService_UpdateCard(t *testing.T) {
	t.Parallel()
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("OK", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:      newListID(),
			BoardID:     newBoardID(),
			Name:        "name",
			Description: "description",
			Position:    1,
		}

		// Create new card.
		card, err := s.CreateCard(&c)
		require.NoError(t, err)

		// Update a single field.
		newName := "new name"
		updatedCard, err := s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)

		// Retrieve card and check.
		other, err := s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, newName, other.Name)
		require.NotNil(t, other.UpdatedAt)
		require.Equal(t, updatedCard, other)

		// Update multiple fields.
		newName = "new name2"
		newDesc := "new description"
		newPosition := float64(2)
		updatedCard, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Name:        &newName,
			Description: &newDesc,
			Position:    &newPosition,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)

		// Retrieve card and check.
		other, err = s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, newName, other.Name)
		require.Equal(t, newDesc, other.Description)
		require.Equal(t, newPosition, other.Position)
		require.Equal(t, updatedCard, other)

		// Set zero values.
		newDesc = ""
		newPosition = 0
		updatedCard, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Description: &newDesc,
			Position:    &newPosition,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)

		// Retrieve card and check.
		other, err = s.Card(card.ID)
		require.NoError(t, err)
		require.Zero(t, other.Description)
		require.Zero(t, other.Position)
		require.Equal(t, updatedCard, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a card that doesn't exist with no patch
		updatedCard, err := s.UpdateCard(newCardID(), &pulpe.CardUpdate{})
		require.Equal(t, pulpe.ErrCardNotFound, err)
		require.Nil(t, updatedCard)

		// Trying to update a card that doesn't exist with a patch
		newName := "new name"
		updatedCard, err = s.UpdateCard(newCardID(), &pulpe.CardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrCardNotFound, err)
		require.Nil(t, updatedCard)
	})
}

func TestCardService_CardsByBoard(t *testing.T) {
	t.Parallel()

	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		boardID1 := newBoardID()
		boardID2 := newBoardID()
		for i := 0; i < 6; i++ {
			c := pulpe.CardCreate{
				ListID:      newListID(),
				Name:        "name",
				Description: "description",
			}

			if i%2 == 0 {
				c.BoardID = boardID1
			} else {
				c.BoardID = boardID2
			}

			// Create new card.
			_, err := s.CreateCard(&c)
			require.NoError(t, err)
		}

		cards, err := s.CardsByBoard(boardID1)
		require.NoError(t, err)
		require.Len(t, cards, 3)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to find cards of a board that doesn't exist.
		cards, err := s.CardsByBoard(newBoardID())
		require.NoError(t, err)
		require.Len(t, cards, 0)
	})
}

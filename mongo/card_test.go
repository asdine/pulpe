package mongo_test

import (
	"fmt"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

// Ensure cards can be created and retrieved.
func TestCardService_CreateCard(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("New", func(t *testing.T) {
		c := pulpe.CardCreate{
			Name:        "YYY",
			ListID:      "ListX",
			BoardID:     "BoardX",
			Description: "MY CARD",
			Position:    1,
		}
		// Create new card.
		card, err := s.CreateCard(&c)
		require.NoError(t, err)
		require.NotZero(t, card.ID)

		// Retrieve card and compare.
		other, err := s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, card, other)
	})

	t.Run("No ListID", func(t *testing.T) {
		// Trying to create a card with no List ID.
		var c pulpe.CardCreate

		_, err := s.CreateCard(&c)
		require.Equal(t, pulpe.ErrCardListIDRequired, err)
	})

	t.Run("No BoardID", func(t *testing.T) {
		// Trying to create a card with no ID.
		c := pulpe.CardCreate{
			ListID: "ListX",
		}

		_, err := s.CreateCard(&c)
		require.Equal(t, pulpe.ErrCardBoardIDRequired, err)
	})
}

// Ensure cards can be retrieved.
func TestCardService_Card(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:  "ListX",
			BoardID: "BoardX",
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
		_, err := s.Card("QQQ")
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCard(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:  "ListX",
			BoardID: "BoardX",
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
		err := s.DeleteCard("QQQ")
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCardsByListID(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		const list1 = pulpe.ListID("List1")
		const list2 = pulpe.ListID("List2")
		const boardID = pulpe.BoardID("board1")

		for i := 0; i < 10; i++ {
			c := pulpe.CardCreate{
				BoardID: boardID,
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
		err := s.DeleteCardsByListID("QQQ")
		require.NoError(t, err)
	})
}

func TestCardService_UpdateCard(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		c := pulpe.CardCreate{
			ListID:      "ListX",
			BoardID:     "BoardX",
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
		newName = ""
		newDesc = ""
		newPosition = 0
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
		require.Zero(t, other.Name)
		require.Zero(t, other.Description)
		require.Zero(t, other.Position)
		require.Equal(t, updatedCard, other)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to update a card that doesn't exist with no patch
		updatedCard, err := s.UpdateCard("QQQ", &pulpe.CardUpdate{})
		require.Equal(t, pulpe.ErrCardNotFound, err)
		require.Nil(t, updatedCard)

		// Trying to update a card that doesn't exist with a patch
		newName := "new name"
		updatedCard, err = s.UpdateCard("QQQ", &pulpe.CardUpdate{Name: &newName})
		require.Equal(t, pulpe.ErrCardNotFound, err)
		require.Nil(t, updatedCard)
	})
}

func TestCardService_CardsByBoard(t *testing.T) {
	session, cleanup := MustGetSession(t)
	defer cleanup()

	s := session.CardService()

	t.Run("Exists", func(t *testing.T) {
		for i := 0; i < 6; i++ {
			c := pulpe.CardCreate{
				ListID:      "ListX",
				BoardID:     pulpe.BoardID(fmt.Sprintf("Board%d", i%2)),
				Name:        "name",
				Description: "description",
			}

			// Create new card.
			_, err := s.CreateCard(&c)
			require.NoError(t, err)
		}

		cards, err := s.CardsByBoard("Board0")
		require.NoError(t, err)
		require.Len(t, cards, 3)
	})

	t.Run("Not found", func(t *testing.T) {
		// Trying to find cards of a board that doesn't exist.
		cards, err := s.CardsByBoard("QQQ")
		require.NoError(t, err)
		require.Len(t, cards, 0)
	})
}

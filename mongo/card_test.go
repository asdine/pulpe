package mongo_test

import (
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

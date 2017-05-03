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
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.CardService()
		l := pulpe.CardCreation{
			Name: "Name",
		}

		// Create new card.
		_, err := s.CreateCard(newListID(), &l)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("New", func(t *testing.T) {
		s := sessions.Red.CardService()

		c := pulpe.CardCreation{
			Name:        "YYY",
			Description: "MY CARD",
			Position:    1,
		}

		// Create new card.
		card, err := s.CreateCard(newList(t, sessions.Red).ID, &c)
		require.NoError(t, err)
		require.NotZero(t, card.ID)
		require.Equal(t, "yyy", card.Slug)

		// Retrieve card and compare.
		other, err := s.Card(card.ID)
		require.NoError(t, err)
		require.Equal(t, card, other)
	})

	t.Run("Slug conflict", func(t *testing.T) {
		s := sessions.Red.CardService()

		listID := newList(t, sessions.Red).ID
		c := pulpe.CardCreation{
			Name:        "ZZZ KK ",
			Description: "MY CARD",
			Position:    1,
		}

		// Create new card.
		card, err := s.CreateCard(listID, &c)
		require.NoError(t, err)
		require.Equal(t, card.Slug, "zzz-kk")

		// Create second card with slightly different name that generates the same slug.
		c = pulpe.CardCreation{
			Name: "  ZZZ   KK ",
		}
		card, err = s.CreateCard(listID, &c)
		require.NoError(t, err)
		require.Equal(t, "zzz-kk-1", card.Slug)
	})
}

// Ensure cards can be retrieved.
func TestCardService_Card(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.CardService()

		// Get a card
		_, err := s.Card("someid")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s1 := sessions.Red.CardService()
		s2 := sessions.Blue.CardService()

		c := pulpe.CardCreation{
			Name: "name1",
		}

		// Create new card as red
		card1, err := s1.CreateCard(newList(t, sessions.Red).ID, &c)
		require.NoError(t, err)

		c.Name = "name2"
		// Create new card as blue
		card2, err := s2.CreateCard(newList(t, sessions.Blue).ID, &c)
		require.NoError(t, err)

		// Retrieve card1 as blue
		other, err := s2.Card(card1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrCardNotFound, err)

		// Retrieve card2 as blue
		other, err = s2.Card(card2.ID)
		require.NoError(t, err)
		require.Equal(t, card2, other)

		// Retrieve card2 as red
		other, err = s1.Card(card2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrCardNotFound, err)

		// Retrieve card1 as red
		other, err = s1.Card(card1.ID)
		require.NoError(t, err)
		require.Equal(t, card1, other)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.CardService()

		// Trying to fetch a card that doesn't exist.
		_, err := s.Card("something")
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.CardService()

		// Trying to delete a card.
		err := s.DeleteCard("something")
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s1 := sessions.Red.CardService()
		s2 := sessions.Blue.CardService()

		c := pulpe.CardCreation{
			Name: "name",
		}

		// Create new card as red
		card1, err := s1.CreateCard(newList(t, sessions.Red).ID, &c)
		require.NoError(t, err)

		// Create new card as blue
		card2, err := s2.CreateCard(newList(t, sessions.Blue).ID, &c)
		require.NoError(t, err)

		// Delete card1 as blue
		err = s2.DeleteCard(card1.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrCardNotFound, err)

		// Delete card2 as red
		err = s1.DeleteCard(card2.ID)
		require.Error(t, err)
		require.Equal(t, pulpe.ErrCardNotFound, err)

		// Delete card1 as red
		err = s1.DeleteCard(card1.ID)
		require.NoError(t, err)

		// Delete card2 as blue
		err = s2.DeleteCard(card2.ID)
		require.NoError(t, err)

		_, err = s1.Card(card1.Slug)
		require.Equal(t, pulpe.ErrCardNotFound, err)

		_, err = s2.Card(card2.Slug)
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.CardService()

		// Trying to delete a card that doesn't exist.
		err := s.DeleteCard(bson.NewObjectId().Hex())
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})
}

func TestCardService_DeleteCardsByListID(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.CardService()

		boardID := newBoard(t, sessions.Red).ID
		list1 := newListWithBoardID(t, sessions.Red, boardID).ID
		list2 := newListWithBoardID(t, sessions.Red, boardID).ID

		var err error
		for i := 0; i < 10; i++ {
			c := pulpe.CardCreation{
				Name: "name",
			}

			if i%2 != 0 {
				_, err = s.CreateCard(list1, &c)
			} else {
				_, err = s.CreateCard(list2, &c)
			}

			require.NoError(t, err)
		}

		// Delete card.
		err = s.DeleteCardsByListID(list1)
		require.NoError(t, err)

		cards, err := s.CardsByBoard(boardID)
		require.NoError(t, err)
		require.Len(t, cards, 5)
		for _, card := range cards {
			require.Equal(t, list2, card.ListID)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Red.CardService()

		l := pulpe.ListCreation{
			Name: "name",
		}

		// Create new List.
		list, err := sessions.Red.ListService().CreateList(newBoard(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Calling with a listID with no associated cards.
		err = s.DeleteCardsByListID(list.ID)
		require.NoError(t, err)
	})
}

func TestCardService_UpdateCard(t *testing.T) {
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Unauthenticated", func(t *testing.T) {
		s := sessions.NoAuth.CardService()

		// Update a card.
		_, err := s.UpdateCard("someid", new(pulpe.CardUpdate))
		require.Error(t, err)
		require.Equal(t, pulpe.ErrUserAuthenticationFailed, err)
	})

	t.Run("OK", func(t *testing.T) {
		s := sessions.Red.CardService()

		c := pulpe.CardCreation{
			Name:        "name",
			Description: "description",
			Position:    1,
		}

		// Create new card.
		card, err := s.CreateCard(newList(t, sessions.Red).ID, &c)
		require.NoError(t, err)

		// Update the name.
		newName := "new name"
		updatedCard, err := s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Name: &newName,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)
		require.Equal(t, newName, updatedCard.Name)
		require.NotNil(t, updatedCard.UpdatedAt)
		require.NotEmpty(t, updatedCard.Slug)

		// Update the name.
		newDesc := "new description"
		updatedCard, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Description: &newDesc,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)
		require.Equal(t, newDesc, updatedCard.Description)
		require.NotNil(t, updatedCard.UpdatedAt)
		require.NotEmpty(t, updatedCard.Slug)

		// Update multiple fields.
		newName = "new name2"
		newDesc = "new description 2"
		newPosition := float64(2)
		updatedCard, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Name:        &newName,
			Description: &newDesc,
			Position:    &newPosition,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)
		require.Equal(t, newName, updatedCard.Name)
		require.Equal(t, newDesc, updatedCard.Description)
		require.Equal(t, newPosition, updatedCard.Position)

		// Set zero values.
		newDesc = ""
		newPosition = 0
		updatedCard, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			Description: &newDesc,
			Position:    &newPosition,
		})
		require.NoError(t, err)
		require.NotNil(t, updatedCard)
		require.Zero(t, updatedCard.Description)
		require.Zero(t, updatedCard.Position)
	})

	t.Run("Change list", func(t *testing.T) {
		s := sessions.Red.CardService()

		c := pulpe.CardCreation{
			Name:        "name",
			Description: "description",
			Position:    1,
		}

		list1 := newList(t, sessions.Red).ID
		list2 := newList(t, sessions.Red).ID
		list3 := newList(t, sessions.Blue).ID

		// Create new card.
		card, err := s.CreateCard(list1, &c)
		require.NoError(t, err)

		// Update the listID for an existing list.
		updatedCard, err := s.UpdateCard(card.ID, &pulpe.CardUpdate{
			ListID: &list2,
		})
		require.NoError(t, err)
		require.Equal(t, list2, updatedCard.ListID)

		// Update the listID for an existing list that's not ours.
		_, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			ListID: &list3,
		})
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)

		// Update the listID for a non existing list.
		listID := "somelist"
		_, err = s.UpdateCard(card.ID, &pulpe.CardUpdate{
			ListID: &listID,
		})
		require.Error(t, err)
		require.Equal(t, pulpe.ErrListNotFound, err)
	})

	t.Run("Bad user", func(t *testing.T) {
		s1 := sessions.Red.CardService()
		s2 := sessions.Blue.CardService()

		l := pulpe.CardCreation{
			Name: "name",
		}

		// Create new Card as red.
		card, err := s1.CreateCard(newList(t, sessions.Red).ID, &l)
		require.NoError(t, err)

		// Update as blue.
		newName := "new name"
		_, err = s2.UpdateCard(card.ID, &pulpe.CardUpdate{
			Name: &newName,
		})
		require.Error(t, err)
		require.Equal(t, pulpe.ErrCardNotFound, err)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.CardService()

		// Trying to update a card that doesn't exist with no patch
		updatedCard, err := s.UpdateCard(newList(t, sessions.Red).ID, &pulpe.CardUpdate{})
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
	sessions, cleanup := MustGetSessions(t)
	defer cleanup()

	t.Run("Exists", func(t *testing.T) {
		s := sessions.Red.CardService()

		boardID1 := newBoard(t, sessions.Red).ID
		boardID2 := newBoard(t, sessions.Red).ID

		l := pulpe.ListCreation{
			Name: "name",
		}

		// Create new List.
		list1, err := sessions.Red.ListService().CreateList(boardID1, &l)
		require.NoError(t, err)

		list2, err := sessions.Red.ListService().CreateList(boardID2, &l)
		require.NoError(t, err)

		for i := 0; i < 6; i++ {
			c := pulpe.CardCreation{
				Name:        "name",
				Description: "description",
			}

			if i%2 == 0 {
				_, err = s.CreateCard(list1.ID, &c)
			} else {
				_, err = s.CreateCard(list2.ID, &c)
			}

			// Create new card.
			require.NoError(t, err)
		}

		cards, err := s.CardsByBoard(boardID1)
		require.NoError(t, err)
		require.Len(t, cards, 3)
	})

	t.Run("Not found", func(t *testing.T) {
		s := sessions.Green.CardService()

		// Trying to find cards of a board that doesn't exist.
		cards, err := s.CardsByBoard(newBoardID())
		require.NoError(t, err)
		require.Len(t, cards, 0)
	})
}

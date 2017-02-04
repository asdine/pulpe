package mongo

import (
	"github.com/blankrobot/pulpe"
)

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// CardService represents a service for managing cards.
type CardService struct {
	session *Session
}

// CreateCard creates a new Card
func (s *CardService) CreateCard(c *pulpe.CardCreate) (*pulpe.Card, error) {
	return nil, nil
}

// Card returns a Card by ID.
func (s *CardService) Card(id pulpe.CardID) (*pulpe.Card, error) {
	return nil, nil
}

// DeleteCard deletes a Card by ID.
func (s *CardService) DeleteCard(id pulpe.CardID) error {
	return nil
}

// UpdateCard updates a Card by ID.
func (s *CardService) UpdateCard(id pulpe.CardID, u *pulpe.CardUpdate) (*pulpe.Card, error) {
	return nil, nil
}

// CardsByBoard returns Cards by board ID.
func (s *CardService) CardsByBoard(boardID pulpe.BoardID) ([]*pulpe.Card, error) {
	return nil, nil
}

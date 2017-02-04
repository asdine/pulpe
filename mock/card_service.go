package mock

import "github.com/blankrobot/pulpe"

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// CardService is a mock service that runs provided functions. Useful for testing.
type CardService struct {
	CreateCardFn      func(card *pulpe.CardCreate) (*pulpe.Card, error)
	CreateCardInvoked bool

	CardFn      func(id pulpe.CardID) (*pulpe.Card, error)
	CardInvoked bool

	DeleteCardFn      func(id pulpe.CardID) error
	DeleteCardInvoked bool

	UpdateCardFn      func(id pulpe.CardID, u *pulpe.CardUpdate) (*pulpe.Card, error)
	UpdateCardInvoked bool

	CardsByBoardFn      func(boardID pulpe.BoardID) ([]*pulpe.Card, error)
	CardsByBoardInvoked bool
}

// CreateCard runs CreateCardFn and sets CreateCardInvoked to true when invoked.
func (s *CardService) CreateCard(card *pulpe.CardCreate) (*pulpe.Card, error) {
	s.CreateCardInvoked = true
	return s.CreateCardFn(card)
}

// Card runs CardFn and sets CardInvoked to true when invoked.
func (s *CardService) Card(id pulpe.CardID) (*pulpe.Card, error) {
	s.CardInvoked = true
	return s.CardFn(id)
}

// DeleteCard runs DeleteCardFn and sets DeleteCardInvoked to true when invoked.
func (s *CardService) DeleteCard(id pulpe.CardID) error {
	s.DeleteCardInvoked = true
	return s.DeleteCardFn(id)
}

// UpdateCard runs UpdateCardFn and sets UpdateCardInvoked to true when invoked.
func (s *CardService) UpdateCard(id pulpe.CardID, u *pulpe.CardUpdate) (*pulpe.Card, error) {
	s.UpdateCardInvoked = true
	return s.UpdateCardFn(id, u)
}

// CardsByBoard runs CardsByBoardFn and sets CardsByBoardInvoked to true when invoked.
func (s *CardService) CardsByBoard(boardID pulpe.BoardID) ([]*pulpe.Card, error) {
	s.CardsByBoardInvoked = true
	return s.CardsByBoardFn(boardID)
}

package mock

import "github.com/blankrobot/pulpe"

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// CardService is a mock service that runs provided functions. Useful for testing.
type CardService struct {
	CreateCardFn      func(card *pulpe.CardCreation) (*pulpe.Card, error)
	CreateCardInvoked bool

	CardFn      func(id string) (*pulpe.Card, error)
	CardInvoked bool

	DeleteCardFn      func(id string) error
	DeleteCardInvoked bool

	DeleteCardsByListIDFn      func(listID string) error
	DeleteCardsByListIDInvoked bool

	DeleteCardsByBoardIDFn      func(boardID string) error
	DeleteCardsByBoardIDInvoked bool

	UpdateCardFn      func(id string, u *pulpe.CardUpdate) (*pulpe.Card, error)
	UpdateCardInvoked bool

	CardsByBoardFn      func(boardID string) ([]*pulpe.Card, error)
	CardsByBoardInvoked bool
}

// CreateCard runs CreateCardFn and sets CreateCardInvoked to true when invoked.
func (s *CardService) CreateCard(card *pulpe.CardCreation) (*pulpe.Card, error) {
	s.CreateCardInvoked = true
	return s.CreateCardFn(card)
}

// Card runs CardFn and sets CardInvoked to true when invoked.
func (s *CardService) Card(id string) (*pulpe.Card, error) {
	s.CardInvoked = true
	return s.CardFn(id)
}

// DeleteCard runs DeleteCardFn and sets DeleteCardInvoked to true when invoked.
func (s *CardService) DeleteCard(id string) error {
	s.DeleteCardInvoked = true
	return s.DeleteCardFn(id)
}

// DeleteCardsByListID runs DeleteCardsByListIDFn and sets DeleteCardsByListIDInvoked to true when invoked.
func (s *CardService) DeleteCardsByListID(listID string) error {
	s.DeleteCardsByListIDInvoked = true
	return s.DeleteCardsByListIDFn(listID)
}

// DeleteCardsByBoardID runs DeleteCardsByBoardIDFn and sets DeleteCardsByBoardIDInvoked to true when invoked.
func (s *CardService) DeleteCardsByBoardID(boardID string) error {
	s.DeleteCardsByBoardIDInvoked = true
	return s.DeleteCardsByBoardIDFn(boardID)
}

// UpdateCard runs UpdateCardFn and sets UpdateCardInvoked to true when invoked.
func (s *CardService) UpdateCard(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
	s.UpdateCardInvoked = true
	return s.UpdateCardFn(id, u)
}

// CardsByBoard runs CardsByBoardFn and sets CardsByBoardInvoked to true when invoked.
func (s *CardService) CardsByBoard(boardID string) ([]*pulpe.Card, error) {
	s.CardsByBoardInvoked = true
	return s.CardsByBoardFn(boardID)
}

package pulpe

import (
	"time"

	shortid "github.com/ventu-io/go-shortid"
)

// CardID represents a Card identifier.
type CardID string

// NewCardID generates a new CardID.
func NewCardID() (CardID, error) {
	id, err := shortid.Generate()
	return CardID(id), err
}

// A Card is a unit of information that is stored in a list.
type Card struct {
	ID          CardID     `json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	ListID      ListID     `json:"listID"`
	BoardID     BoardID    `json:"boardID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Position    float64    `json:"position"`
}

// CardCreate is used to create a Card.
type CardCreate struct {
	ListID      ListID  `json:"listID"`
	BoardID     BoardID `json:"boardID"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Position    float64 `json:"position"`
}

// CardUpdate is used to update a Card.
type CardUpdate struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Position    *float64 `json:"position"`
}

// CardService represents a service for managing cards.
type CardService interface {
	CreateCard(card *CardCreate) (*Card, error)
	Card(id CardID) (*Card, error)
	DeleteCard(id CardID) error
	UpdateCard(id CardID, u *CardUpdate) (*Card, error)
	CardsByBoard(boardID BoardID) ([]*Card, error)
}

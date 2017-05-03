package pulpe

import "time"

// Card errors
const (
	ErrCardNotFound = Error("card not found")
)

// A Card is a unit of information that is stored in a list.
type Card struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	OwnerID     string     `json:"ownerID"`
	ListID      string     `json:"listID"`
	BoardID     string     `json:"boardID"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Position    float64    `json:"position"`
}

// CardCreation is used to create a Card.
type CardCreation struct {
	Name        string
	Description string
	Position    float64
}

// CardUpdate is used to update a Card.
type CardUpdate struct {
	Name        *string
	Description *string
	Position    *float64
	ListID      *string
}

// CardService represents a service for managing cards.
type CardService interface {
	CreateCard(listID string, card *CardCreation) (*Card, error)
	Card(id string) (*Card, error)
	DeleteCard(id string) error
	UpdateCard(id string, u *CardUpdate) (*Card, error)
	DeleteCardsByListID(listID string) error
	DeleteCardsByBoardID(boardID string) error
	CardsByBoard(boardID string) ([]*Card, error)
}

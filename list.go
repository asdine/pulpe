package pulpe

import (
	"time"

	shortid "github.com/ventu-io/go-shortid"
)

// ListID represents a List identifier.
type ListID string

// NewListID generates a new ListID.
func NewListID() (ListID, error) {
	id, err := shortid.Generate()
	return ListID(id), err
}

// A List is a container of cards.
type List struct {
	ID        ListID     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	BoardID   BoardID    `json:"boardID"`
	Name      string     `json:"name"`
}

// ListCreate is used to create a List.
type ListCreate struct {
	BoardID BoardID `json:"boardID"`
	Name    string  `json:"name"`
}

// ListUpdate is used to update a List.
type ListUpdate struct {
	Name *string `json:"name"`
}

// ListService represents a service for managing lists.
type ListService interface {
	CreateList(list *ListCreate) (*List, error)
	List(id ListID) (*List, error)
	DeleteList(id ListID) error
	UpdateList(id ListID, u *ListUpdate) (*List, error)
	ListsByBoard(boardID BoardID) ([]*List, error)
}

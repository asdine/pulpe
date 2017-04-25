package pulpe

import "time"

// List errors
const (
	ErrListNotFound = Error("list not found")
)

// A List is a container of cards.
type List struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	OwnerID   string     `json:"ownerID"`
	BoardID   string     `json:"boardID"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
}

// ListCreation is used to create a List.
type ListCreation struct {
	Name string
}

// ListUpdate is used to update a List.
type ListUpdate struct {
	Name *string
}

// ListService represents a service for managing lists.
type ListService interface {
	CreateList(boardID string, list *ListCreation) (*List, error)
	List(id string) (*List, error)
	DeleteList(id string) error
	UpdateList(id string, u *ListUpdate) (*List, error)
	DeleteListsByBoardID(boardID string) error
	ListsByBoard(boardID string) ([]*List, error)
}

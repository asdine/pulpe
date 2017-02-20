package pulpe

import "time"

// A List is a container of cards.
type List struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	BoardID   string     `json:"boardID"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
}

// ListCreate is used to create a List.
type ListCreate struct {
	BoardID string
	Name    string
}

// ListUpdate is used to update a List.
type ListUpdate struct {
	Name *string
}

// ListService represents a service for managing lists.
type ListService interface {
	CreateList(list *ListCreate) (*List, error)
	List(id string) (*List, error)
	DeleteList(id string) error
	DeleteListsByBoardID(boardID string) error
	UpdateList(id string, u *ListUpdate) (*List, error)
	ListsByBoard(boardID string) ([]*List, error)
}

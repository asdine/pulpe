package pulpe

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// ListID represents a List identifier.
type ListID string

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

// Validate list creation payload.
func (l *ListCreate) Validate() error {
	l.Name = strings.TrimSpace(l.Name)
	return validation.ValidateStruct(l,
		validation.Field(&l.BoardID, validation.Required, is.Alphanumeric, validation.Length(1, 64)),
		validation.Field(&l.Name, validation.Required, validation.Length(1, 64)),
	)
}

// ListUpdate is used to update a List.
type ListUpdate struct {
	Name *string `json:"name"`
}

// Validate list update payload.
func (l *ListUpdate) Validate() error {
	if l.Name == nil {
		return nil
	}

	name := strings.TrimSpace(*l.Name)

	return validation.Errors{
		"name": validation.Validate(name, validation.Required, validation.Length(1, 64)),
	}.Filter()
}

// ListService represents a service for managing lists.
type ListService interface {
	CreateList(list *ListCreate) (*List, error)
	List(id ListID) (*List, error)
	DeleteList(id ListID) error
	DeleteListsByBoardID(boardID BoardID) error
	UpdateList(id ListID, u *ListUpdate) (*List, error)
	ListsByBoard(boardID BoardID) ([]*List, error)
}

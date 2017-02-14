package pulpe

import (
	"strings"
	"time"

	"github.com/blankrobot/pulpe/validation"
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
	BoardID BoardID `json:"boardID" valid:"required,stringlength(1|64),alphanum"`
	Name    string  `json:"name" valid:"required,stringlength(1|64)"`
}

// Validate list creation payload.
func (l *ListCreate) Validate(s Session) error {
	verr := validation.Validate(l)

	// validate boardID existence if boardID is valid.
	if validation.LastError(verr, "boardID") == nil {
		_, err := s.BoardService().Board(l.BoardID)
		if err != nil && err != ErrBoardNotFound {
			return err
		}
		if err == ErrBoardNotFound {
			verr = validation.AddError(verr, "boardID", err)
		}
	}

	return verr
}

// ListUpdate is used to update a List.
type ListUpdate struct {
	Name *string `json:"name" valid:"required,stringlength(1|64)"`
}

// Validate list update payload.
func (l *ListUpdate) Validate() error {
	if l.Name == nil {
		return nil
	}

	*l.Name = strings.TrimSpace(*l.Name)

	return validation.Validate(l)
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

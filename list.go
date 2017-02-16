package pulpe

import (
	"strings"
	"time"

	"github.com/blankrobot/pulpe/validation"
)

// A List is a container of cards.
type List struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	BoardID   string     `json:"boardID"`
	Name      string     `json:"name"`
}

// ListCreate is used to create a List.
type ListCreate struct {
	BoardID string `json:"boardID" valid:"required,stringlength(1|64),alphanum"`
	Name    string `json:"name" valid:"required,stringlength(1|64)"`
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
	List(id string) (*List, error)
	DeleteList(id string) error
	DeleteListsByBoardID(boardID string) error
	UpdateList(id string, u *ListUpdate) (*List, error)
	ListsByBoard(boardID string) ([]*List, error)
}

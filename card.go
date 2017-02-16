package pulpe

import (
	"errors"
	"strings"
	"time"

	"github.com/blankrobot/pulpe/validation"
)

// A Card is a unit of information that is stored in a list.
type Card struct {
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	ListID      string     `json:"listID"`
	BoardID     string     `json:"boardID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Position    float64    `json:"position"`
}

// CardCreate is used to create a Card.
type CardCreate struct {
	ListID      string  `json:"listID" valid:"required,stringlength(1|64),alphanum"`
	BoardID     string  `json:"boardID" valid:"required,stringlength(1|64),alphanum"`
	Name        string  `json:"name" valid:"required,stringlength(1|64)"`
	Description string  `json:"description" valid:"stringlength(1|100000)"`
	Position    float64 `json:"position"`
}

// Validate list creation payload.
func (c *CardCreate) Validate(s Session) error {
	verr := validation.Validate(c)

	// validate boardID existence if boardID is valid.
	if validation.LastError(verr, "boardID") == nil {
		_, err := s.BoardService().Board(c.BoardID)
		if err != nil && err != ErrBoardNotFound {
			return err
		}
		if err == ErrBoardNotFound {
			verr = validation.AddError(verr, "boardID", err)
		}
	}

	// validate listID existence if listID is valid.
	if validation.LastError(verr, "listID") == nil {
		_, err := s.ListService().List(c.ListID)
		if err != nil && err != ErrListNotFound {
			return err
		}
		if err == ErrListNotFound {
			verr = validation.AddError(verr, "listID", err)
		}
	}

	// validate position
	if c.Position < 0 {
		verr = validation.AddError(verr, "position", errors.New("position should be greater than zero"))
	}

	return verr
}

// CardUpdate is used to update a Card.
type CardUpdate struct {
	Name        *string  `json:"name" valid:"stringlength(1|64)"`
	Description *string  `json:"description" valid:"stringlength(1|100000)"`
	Position    *float64 `json:"position"`
}

// Validate card update payload.
func (c *CardUpdate) Validate() error {
	if c.Name != nil {
		*c.Name = strings.TrimSpace(*c.Name)
	}

	if c.Description != nil {
		*c.Description = strings.TrimSpace(*c.Description)
	}

	err := validation.Validate(c)
	if c.Name != nil && *c.Name == "" {
		err = validation.AddError(err, "name", errors.New("name should not be empty"))
	}

	if c.Position != nil && *c.Position < 0 {
		err = validation.AddError(err, "position", errors.New("position should be greater than zero"))
	}

	return err
}

// CardService represents a service for managing cards.
type CardService interface {
	CreateCard(card *CardCreate) (*Card, error)
	Card(id string) (*Card, error)
	DeleteCard(id string) error
	DeleteCardsByListID(listID string) error
	DeleteCardsByBoardID(boardID string) error
	UpdateCard(id string, u *CardUpdate) (*Card, error)
	CardsByBoard(boardID string) ([]*Card, error)
}

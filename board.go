package pulpe

import (
	"encoding/json"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// BoardID represents a Board identifier.
type BoardID string

// A Board is a container of lists.
type Board struct {
	ID        BoardID          `json:"id"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt *time.Time       `json:"updatedAt,omitempty"`
	Name      string           `json:"name"`
	Slug      string           `json:"slug"`
	Lists     []*List          `json:"lists"`
	Cards     []*Card          `json:"cards"`
	Settings  *json.RawMessage `json:"settings,omitempty"`
}

// BoardCreate is used to create a board.
type BoardCreate struct {
	Name     string           `json:"name"`
	Settings *json.RawMessage `json:"settings"`
}

// Validate board creation payload.
func (b *BoardCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required, validation.Length(1, 64)),
	)
}

// BoardUpdate is used to update a board.
type BoardUpdate struct {
	Name     *string          `json:"name"`
	Settings *json.RawMessage `json:"settings"`
}

// Validate board update payload.
func (b *BoardUpdate) Validate() error {
	if b.Name == nil {
		return nil
	}

	name := strings.TrimSpace(*b.Name)

	return validation.Errors{
		"name": validation.Validate(name, validation.Required, validation.Length(1, 64)),
	}.Filter()
}

// BoardService represents a service for managing boards.
type BoardService interface {
	CreateBoard(board *BoardCreate) (*Board, error)
	Board(id BoardID) (*Board, error)
	Boards() ([]*Board, error)
	DeleteBoard(id BoardID) error
	UpdateBoard(id BoardID, u *BoardUpdate) (*Board, error)
}

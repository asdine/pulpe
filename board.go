package pulpe

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/blankrobot/pulpe/validation"
)

// A Board is a container of lists.
type Board struct {
	ID        string           `json:"id"`
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
	Name     string           `json:"name" valid:"required,stringlength(1|64)"`
	Settings *json.RawMessage `json:"settings" valid:"pulpe-json"`
}

// Validate board creation payload.
func (b *BoardCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	return validation.Validate(b)
}

// BoardUpdate is used to update a board.
type BoardUpdate struct {
	Name     *string          `json:"name" valid:"stringlength(1|64)"`
	Settings *json.RawMessage `json:"settings" valid:"pulpe-json"`
}

// Validate board update payload.
func (b *BoardUpdate) Validate() error {
	if b.Name != nil {
		*b.Name = strings.TrimSpace(*b.Name)
	}

	err := validation.Validate(b)
	if b.Name != nil && *b.Name == "" {
		err = validation.AddError(err, "name", errors.New("name should not be empty"))
	}

	return err
}

// BoardService represents a service for managing boards.
type BoardService interface {
	CreateBoard(board *BoardCreate) (*Board, error)
	Board(id string) (*Board, error)
	Boards(map[string]string) ([]*Board, error)
	DeleteBoard(id string) error
	UpdateBoard(id string, u *BoardUpdate) (*Board, error)
}

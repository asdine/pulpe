package pulpe

import (
	"encoding/json"
	"time"
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

// BoardUpdate is used to update a board.
type BoardUpdate struct {
	Name     *string          `json:"name"`
	Settings *json.RawMessage `json:"settings"`
}

// BoardService represents a service for managing boards.
type BoardService interface {
	CreateBoard(board *BoardCreate) (*Board, error)
	Board(id BoardID) (*Board, error)
	Boards() ([]*Board, error)
	DeleteBoard(id BoardID) error
	UpdateBoard(id BoardID, u *BoardUpdate) (*Board, error)
}

package pulpe

import "time"

// A Board is a container of lists.
type Board struct {
	ID        string     `json:"id"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Name      string     `json:"name"`
	Lists     []*List    `json:"lists,omitempty"`
	Cards     []*Card    `json:"cards,omitempty"`
}

// BoardCreate is used to create a board.
type BoardCreate struct {
	Name string
}

// BoardUpdate is used to update a board.
type BoardUpdate struct {
	Name *string
}

// BoardService represents a service for managing boards.
type BoardService interface {
	CreateBoard(board *BoardCreate) (*Board, error)
	Board(id string) (*Board, error)
	Boards() ([]*Board, error)
	DeleteBoard(id string) error
	UpdateBoard(id string, u *BoardUpdate) (*Board, error)
}

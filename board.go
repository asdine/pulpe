package pulpe

import "time"

// Board errors
const (
	ErrBoardNotFound = Error("board not found")
)

// A Board is a container of lists.
type Board struct {
	ID        string     `json:"id"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	OwnerID   string     `json:"ownerID"`
	Name      string     `json:"name"`
	Lists     []*List    `json:"lists,omitempty"`
	Cards     []*Card    `json:"cards,omitempty"`
}

// BoardCreation is used to create a board.
type BoardCreation struct {
	Name string
}

// BoardUpdate is used to update a board.
type BoardUpdate struct {
	Name *string
}

// BoardService represents a service for managing boards.
type BoardService interface {
	CreateBoard(board *BoardCreation) (*Board, error)
	Board(id string, options ...BoardGetOption) (*Board, error)
	BoardByOwnerAndSlug(owner, slug string, options ...BoardGetOption) (*Board, error)
	Boards() ([]*Board, error)
	DeleteBoard(id string) error
	UpdateBoard(id string, u *BoardUpdate) (*Board, error)
}

// BoardGetOption is a function used to customize the way a board is fetched.
type BoardGetOption func(*BoardGetOptions)

// BoardGetOptions contains the list of options to customize the way a board is fetched.
type BoardGetOptions struct {
	WithCards bool
	WithLists bool
}

// WithCards is used to tell the BoardService to also fetch cards.
func WithCards() BoardGetOption {
	return func(b *BoardGetOptions) {
		b.WithCards = true
	}
}

// WithLists is used to tell the BoardService to also fetch lists.
func WithLists() BoardGetOption {
	return func(b *BoardGetOptions) {
		b.WithCards = true
	}
}

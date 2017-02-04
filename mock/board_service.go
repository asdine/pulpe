package mock

import "github.com/blankrobot/pulpe"

// Ensure BoardService implements pulpe.BoardService.
var _ pulpe.BoardService = new(BoardService)

// BoardService is a mock service that runs provided functions. Useful for testing.
type BoardService struct {
	CreateBoardFn      func(Board *pulpe.BoardCreate) (*pulpe.Board, error)
	CreateBoardInvoked bool

	BoardFn      func(id pulpe.BoardID) (*pulpe.Board, error)
	BoardInvoked bool

	BoardsFn      func() ([]*pulpe.Board, error)
	BoardsInvoked bool

	DeleteBoardFn      func(id pulpe.BoardID) error
	DeleteBoardInvoked bool

	UpdateBoardFn      func(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error)
	UpdateBoardInvoked bool
}

// CreateBoard runs CreateBoardFn and sets CreateBoardInvoked to true when invoked.
func (s *BoardService) CreateBoard(Board *pulpe.BoardCreate) (*pulpe.Board, error) {
	s.CreateBoardInvoked = true
	return s.CreateBoardFn(Board)
}

// Board runs BoardFn and sets BoardInvoked to true when invoked.
func (s *BoardService) Board(id pulpe.BoardID) (*pulpe.Board, error) {
	s.BoardInvoked = true
	return s.BoardFn(id)
}

// Boards runs BoardsFn and sets BoardsInvoked to true when invoked.
func (s *BoardService) Boards() ([]*pulpe.Board, error) {
	s.BoardsInvoked = true
	return s.BoardsFn()
}

// DeleteBoard runs DeleteBoardFn and sets DeleteBoardInvoked to true when invoked.
func (s *BoardService) DeleteBoard(id pulpe.BoardID) error {
	s.DeleteBoardInvoked = true
	return s.DeleteBoardFn(id)
}

// UpdateBoard runs UpdateBoardFn and sets UpdateBoardInvoked to true when invoked.
func (s *BoardService) UpdateBoard(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
	s.UpdateBoardInvoked = true
	return s.UpdateBoardFn(id, u)
}

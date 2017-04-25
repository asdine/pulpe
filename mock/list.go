package mock

import "github.com/blankrobot/pulpe"

// Ensure ListService implements pulpe.ListService.
var _ pulpe.ListService = new(ListService)

// ListService is a mock service that runs provided functions. Useful for testing.
type ListService struct {
	CreateListFn      func(boardID string, list *pulpe.ListCreation) (*pulpe.List, error)
	CreateListInvoked bool

	ListFn      func(id string) (*pulpe.List, error)
	ListInvoked bool

	DeleteListFn      func(id string) error
	DeleteListInvoked bool

	DeleteListsByBoardIDFn      func(boardID string) error
	DeleteListsByBoardIDInvoked bool

	UpdateListFn      func(id string, u *pulpe.ListUpdate) (*pulpe.List, error)
	UpdateListInvoked bool

	ListsByBoardFn      func(boardID string) ([]*pulpe.List, error)
	ListsByBoardInvoked bool
}

// CreateList runs CreateListFn and sets CreateListInvoked to true when invoked.
func (s *ListService) CreateList(boardID string, list *pulpe.ListCreation) (*pulpe.List, error) {
	s.CreateListInvoked = true
	return s.CreateListFn(boardID, list)
}

// List runs ListFn and sets ListInvoked to true when invoked.
func (s *ListService) List(id string) (*pulpe.List, error) {
	s.ListInvoked = true
	return s.ListFn(id)
}

// DeleteList runs DeleteListFn and sets DeleteListInvoked to true when invoked.
func (s *ListService) DeleteList(id string) error {
	s.DeleteListInvoked = true
	return s.DeleteListFn(id)
}

// DeleteListsByBoardID runs DeleteListsByBoardIDFn and sets DeleteListsByBoardIDInvoked to true when invoked.
func (s *ListService) DeleteListsByBoardID(boardID string) error {
	s.DeleteListsByBoardIDInvoked = true
	return s.DeleteListsByBoardIDFn(boardID)
}

// UpdateList runs UpdateListFn and sets UpdateListInvoked to true when invoked.
func (s *ListService) UpdateList(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
	s.UpdateListInvoked = true
	return s.UpdateListFn(id, u)
}

// ListsByBoard runs ListsByBoardFn and sets ListsByBoardInvoked to true when invoked.
func (s *ListService) ListsByBoard(boardID string) ([]*pulpe.List, error) {
	s.ListsByBoardInvoked = true
	return s.ListsByBoardFn(boardID)
}

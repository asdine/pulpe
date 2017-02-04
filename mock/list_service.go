package mock

import "github.com/blankrobot/pulpe"

// Ensure ListService implements pulpe.ListService.
var _ pulpe.ListService = new(ListService)

// ListService is a mock service that runs provided functions. Useful for testing.
type ListService struct {
	CreateListFn      func(list *pulpe.ListCreate) (*pulpe.List, error)
	CreateListInvoked bool

	ListFn      func(id pulpe.ListID) (*pulpe.List, error)
	ListInvoked bool

	DeleteListFn      func(id pulpe.ListID) error
	DeleteListInvoked bool

	UpdateListFn      func(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error)
	UpdateListInvoked bool

	ListsByBoardFn      func(boardID pulpe.BoardID) ([]*pulpe.List, error)
	ListsByBoardInvoked bool
}

// CreateList runs CreateListFn and sets CreateListInvoked to true when invoked.
func (s *ListService) CreateList(list *pulpe.ListCreate) (*pulpe.List, error) {
	s.CreateListInvoked = true
	return s.CreateListFn(list)
}

// List runs ListFn and sets ListInvoked to true when invoked.
func (s *ListService) List(id pulpe.ListID) (*pulpe.List, error) {
	s.ListInvoked = true
	return s.ListFn(id)
}

// DeleteList runs DeleteListFn and sets DeleteListInvoked to true when invoked.
func (s *ListService) DeleteList(id pulpe.ListID) error {
	s.DeleteListInvoked = true
	return s.DeleteListFn(id)
}

// UpdateList runs UpdateListFn and sets UpdateListInvoked to true when invoked.
func (s *ListService) UpdateList(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
	s.UpdateListInvoked = true
	return s.UpdateListFn(id, u)
}

// ListsByBoard runs ListsByBoardFn and sets ListsByBoardInvoked to true when invoked.
func (s *ListService) ListsByBoard(boardID pulpe.BoardID) ([]*pulpe.List, error) {
	s.ListsByBoardInvoked = true
	return s.ListsByBoardFn(boardID)
}

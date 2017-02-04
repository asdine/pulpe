package mongo

import "github.com/blankrobot/pulpe"

// Ensure ListService implements pulpe.ListService.
var _ pulpe.ListService = new(ListService)

// ListService represents a service for managing lists.
type ListService struct {
	session *Session
}

// CreateList creates a new List
func (s *ListService) CreateList(l *pulpe.ListCreate) (*pulpe.List, error) {
	return nil, nil
}

// List returns a List by ID.
func (s *ListService) List(id pulpe.ListID) (*pulpe.List, error) {
	return nil, nil
}

// DeleteList deletes a List by ID.
func (s *ListService) DeleteList(id pulpe.ListID) error {
	return nil
}

// UpdateList updates a List by ID.
func (s *ListService) UpdateList(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
	return nil, nil
}

// ListsByBoard returns all the lists of a given board.
func (s *ListService) ListsByBoard(boardID pulpe.BoardID) ([]*pulpe.List, error) {
	return nil, nil
}

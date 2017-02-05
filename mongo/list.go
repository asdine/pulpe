package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
)

const listCol = "lists"

// Ensure ListService implements pulpe.ListService.
var _ pulpe.ListService = new(ListService)

// List representation stored in MongoDB.
type List struct {
	ID        bson.ObjectId `bson:"_id"`
	PublicID  string        `bson:"publicID"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	BoardID   string        `bson:"boardID"`
	Name      string        `bson:"name"`
}

// ToMongoList creates a mongo list from a pulpe list.
func ToMongoList(p *pulpe.List) *List {
	return &List{
		ID:        bson.NewObjectId(),
		PublicID:  string(p.ID),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		BoardID:   string(p.BoardID),
		Name:      p.Name,
	}
}

// FromMongoList creates a pulpe list from a mongo list.
func FromMongoList(l *List) *pulpe.List {
	p := pulpe.List{
		ID:        pulpe.ListID(l.PublicID),
		CreatedAt: l.CreatedAt,
		BoardID:   pulpe.BoardID(l.BoardID),
		Name:      l.Name,
	}

	if l.UpdatedAt != nil {
		t := (*l.UpdatedAt).UTC()
		p.UpdatedAt = &t
	}

	return &p
}

// ListService represents a service for managing lists.
type ListService struct {
	session *Session
}

// CreateList creates a new List
func (s *ListService) CreateList(l *pulpe.ListCreate) (*pulpe.List, error) {
	if l.BoardID == "" {
		return nil, pulpe.ErrListBoardIDRequired
	}

	var err error
	list := pulpe.List{
		BoardID: l.BoardID,
		Name:    l.Name,
	}

	list.ID, err = pulpe.NewListID()
	if err != nil {
		return nil, err
	}

	list.CreatedAt = s.session.now

	return &list, s.session.db.C(listCol).Insert(ToMongoList(&list))
}

// List returns a List by ID.
func (s *ListService) List(id pulpe.ListID) (*pulpe.List, error) {
	var l List

	err := s.session.db.C(listCol).Find(bson.M{"publicID": string(id)}).One(&l)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrListNotFound
		}

		return nil, err
	}

	return FromMongoList(&l), nil
}

// DeleteList deletes a List by ID.
func (s *ListService) DeleteList(id pulpe.ListID) error {
	err := s.session.db.C(listCol).Remove(bson.M{"publicID": string(id)})
	if err == mgo.ErrNotFound {
		return pulpe.ErrListNotFound
	}

	return err
}

// UpdateList updates a List by ID.
func (s *ListService) UpdateList(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
	return nil, nil
}

// ListsByBoard returns all the lists of a given board.
func (s *ListService) ListsByBoard(boardID pulpe.BoardID) ([]*pulpe.List, error) {
	return nil, nil
}

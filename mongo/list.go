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
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	BoardID   bson.ObjectId `bson:"boardID"`
	Name      string        `bson:"name"`
}

// ToMongoList creates a mongo list from a pulpe list.
func ToMongoList(p *pulpe.List) *List {
	var id bson.ObjectId

	if p.ID == "" {
		id = bson.NewObjectId()
		p.ID = pulpe.ListID(id.Hex())
	} else {
		id = bson.ObjectIdHex(string(p.ID))
	}

	return &List{
		ID:        id,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		BoardID:   bson.ObjectIdHex(string(p.BoardID)),
		Name:      p.Name,
	}
}

// FromMongoList creates a pulpe list from a mongo list.
func FromMongoList(l *List) *pulpe.List {
	p := pulpe.List{
		ID:        pulpe.ListID(l.ID.Hex()),
		CreatedAt: l.CreatedAt,
		BoardID:   pulpe.BoardID(l.BoardID.Hex()),
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

func (s *ListService) ensureIndexes() error {
	col := s.session.db.C(listCol)

	// boardID
	index := mgo.Index{
		Key:    []string{"boardID"},
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateList creates a new List
func (s *ListService) CreateList(l *pulpe.ListCreate) (*pulpe.List, error) {
	if l.BoardID == "" {
		return nil, pulpe.ErrListBoardIDRequired
	}

	list := pulpe.List{
		CreatedAt: s.session.now,
		BoardID:   l.BoardID,
		Name:      l.Name,
	}

	return &list, s.session.db.C(listCol).Insert(ToMongoList(&list))
}

// List returns a List by ID.
func (s *ListService) List(id pulpe.ListID) (*pulpe.List, error) {
	var l List

	err := s.session.db.C(listCol).FindId(bson.ObjectIdHex(string(id))).One(&l)
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
	err := s.session.db.C(listCol).RemoveId(bson.ObjectIdHex(string(id)))
	if err == mgo.ErrNotFound {
		return pulpe.ErrListNotFound
	}

	return err
}

// DeleteListsByBoardID deletes all the lists of a board.
func (s *ListService) DeleteListsByBoardID(boardID pulpe.BoardID) error {
	_, err := s.session.db.C(listCol).RemoveAll(bson.M{"boardID": bson.ObjectIdHex(string(boardID))})
	return err
}

// UpdateList updates a List by ID.
func (s *ListService) UpdateList(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
	col := s.session.db.C(listCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
	}

	if len(patch) == 0 {
		return s.List(id)
	}

	err := col.UpdateId(
		bson.ObjectIdHex(string(id)),
		bson.M{
			"$set":         patch,
			"$currentDate": bson.M{"updatedAt": true},
		})
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrListNotFound
		}

		return nil, err
	}

	return s.List(id)
}

// ListsByBoard returns all the lists of a given board.
func (s *ListService) ListsByBoard(boardID pulpe.BoardID) ([]*pulpe.List, error) {
	var lists []List

	// TODO set a limit
	err := s.session.db.C(listCol).Find(bson.M{"boardID": bson.ObjectIdHex(string(boardID))}).All(&lists)
	if err != nil {
		return nil, err
	}

	list := make([]*pulpe.List, len(lists))
	for i := range lists {
		list[i] = FromMongoList(&lists[i])
	}

	return list, nil
}

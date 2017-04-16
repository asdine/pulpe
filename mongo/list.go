package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Machiel/slugify"
	"github.com/blankrobot/pulpe"
)

const listCol = "lists"

// Ensure ListService implements pulpe.ListService.
var _ pulpe.ListService = new(ListService)

// List representation stored in MongoDB.
type List struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	BoardID   bson.ObjectId `bson:"boardID"`
	Name      string        `bson:"name"`
	Slug      string        `bson:"slug"`
}

// ToMongoList creates a mongo list from a pulpe list.
func ToMongoList(p *pulpe.List) *List {
	id := bson.NewObjectId()
	p.ID = id.Hex()
	p.CreatedAt = id.Time()

	return &List{
		ID:        id,
		UpdatedAt: p.UpdatedAt,
		BoardID:   bson.ObjectIdHex(p.BoardID),
		Name:      p.Name,
		Slug:      p.Slug,
	}
}

// FromMongoList creates a pulpe list from a mongo list.
func FromMongoList(l *List) *pulpe.List {
	p := pulpe.List{
		ID:        l.ID.Hex(),
		CreatedAt: l.ID.Time(),
		BoardID:   l.BoardID.Hex(),
		Name:      l.Name,
		Slug:      l.Slug,
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
		Key:    []string{"boardID", "slug"},
		Unique: true,
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateList creates a new List
func (s *ListService) CreateList(lc *pulpe.ListCreation) (*pulpe.List, error) {
	var err error

	// generate slug
	slug := slugify.Slugify(lc.Name)

	list := pulpe.List{
		BoardID: lc.BoardID,
		Name:    lc.Name,
		Slug:    slug,
	}

	l := ToMongoList(&list)
	col := s.session.db.C(listCol)

	list.Slug, err = resolveSlugAndDo(col, newListRecorder(l), func(rec recorder) error {
		return col.Insert(rec.elem())
	})

	return &list, err
}

// List returns a List by ID.
func (s *ListService) List(id string) (*pulpe.List, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrListNotFound
	}

	var l List

	err := s.session.db.C(listCol).FindId(bson.ObjectIdHex(id)).One(&l)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrListNotFound
		}

		return nil, err
	}

	return FromMongoList(&l), nil
}

// DeleteList deletes a List by ID.
func (s *ListService) DeleteList(id string) error {
	if !bson.IsObjectIdHex(id) {
		return pulpe.ErrListNotFound
	}

	err := s.session.db.C(listCol).RemoveId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return pulpe.ErrListNotFound
	}

	return err
}

// DeleteListsByBoardID deletes all the lists of a board.
func (s *ListService) DeleteListsByBoardID(boardID string) error {
	if !bson.IsObjectIdHex(boardID) {
		return pulpe.ErrBoardNotFound
	}

	_, err := s.session.db.C(listCol).RemoveAll(bson.M{"boardID": bson.ObjectIdHex(boardID)})
	return err
}

// UpdateList updates a List by ID.
func (s *ListService) UpdateList(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrListNotFound
	}

	var err error
	var l List

	col := s.session.db.C(listCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		l.Slug = slugify.Slugify(*u.Name)
	}

	if len(patch) == 0 {
		return s.List(id)
	}

	l.Slug, err = resolveSlugAndDo(col, newListRecorder(&l), func(rec recorder) error {
		slug := rec.getSlug()
		if slug != "" {
			patch["slug"] = slug
		}

		return col.UpdateId(
			bson.ObjectIdHex(id),
			bson.M{
				"$set":         patch,
				"$currentDate": bson.M{"updatedAt": true},
			})
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
func (s *ListService) ListsByBoard(boardID string) ([]*pulpe.List, error) {
	if !bson.IsObjectIdHex(boardID) {
		return nil, pulpe.ErrBoardNotFound
	}

	var lists []List

	// TODO set a limit
	err := s.session.db.C(listCol).Find(bson.M{"boardID": bson.ObjectIdHex(boardID)}).Sort("_id").All(&lists)
	if err != nil {
		return nil, err
	}

	list := make([]*pulpe.List, len(lists))
	for i := range lists {
		list[i] = FromMongoList(&lists[i])
	}

	return list, nil
}

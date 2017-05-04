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

// list representation stored in MongoDB.
type list struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	OwnerID   string        `bson:"ownerID"`
	BoardID   string        `bson:"boardID"`
	Name      string        `bson:"name"`
	Slug      string        `bson:"slug"`
	Position  float64       `bson:"position"`
}

// toPulpeList creates a pulpe list from a mongo list.
func (l *list) toPulpeList() *pulpe.List {
	p := pulpe.List{
		ID:        l.ID.Hex(),
		CreatedAt: l.ID.Time().UTC(),
		OwnerID:   l.OwnerID,
		BoardID:   l.BoardID,
		Name:      l.Name,
		Slug:      l.Slug,
		Position:  l.Position,
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
	store   listStore
}

func (s *ListService) ensureIndexes() error {
	col := s.session.db.C(listCol)

	// boardID and slug
	index := mgo.Index{
		Key:    []string{"boardID", "slug"},
		Unique: true,
		Sparse: true,
	}

	err := col.EnsureIndex(index)
	if err != nil {
		return err
	}

	index = mgo.Index{
		Key:    []string{"_id", "ownerID"},
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateList creates a new List
func (s *ListService) CreateList(boardID string, lc *pulpe.ListCreation) (*pulpe.List, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	board, err := s.session.BoardService().Board(boardID)
	if err != nil {
		return nil, err
	}

	if board.Owner.ID != user.ID {
		return nil, pulpe.ErrBoardNotFound
	}

	l := list{
		ID:       bson.NewObjectId(),
		OwnerID:  user.ID,
		BoardID:  board.ID,
		Name:     lc.Name,
		Slug:     slugify.Slugify(lc.Name),
		Position: lc.Position,
	}

	err = s.store.createList(&l)
	if err != nil {
		return nil, err
	}

	return l.toPulpeList(), err
}

// List returns a List.
func (s *ListService) List(id string) (*pulpe.List, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	l, err := s.store.listByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrListNotFound
		}

		return nil, err
	}

	return l.toPulpeList(), nil
}

// DeleteList deletes a List.
func (s *ListService) DeleteList(id string) error {
	user, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	if !bson.IsObjectIdHex(id) {
		return pulpe.ErrListNotFound
	}

	err = s.store.deleteListByOwnerIDAndID(user.ID, bson.ObjectIdHex(id))
	if err != nil {
		if err == mgo.ErrNotFound {
			return pulpe.ErrListNotFound

		}

		return err
	}

	return s.session.CardService().DeleteCardsByListID(id)
}

// DeleteListsByBoardID deletes all the lists of a board.
func (s *ListService) DeleteListsByBoardID(boardID string) error {
	return s.store.deleteListsByBoardID(boardID)
}

// UpdateList updates a List by ID.
func (s *ListService) UpdateList(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrListNotFound
	}

	var newSlug string

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		newSlug = slugify.Slugify(*u.Name)
	}

	if u.Position != nil {
		patch["position"] = *u.Position
	}

	if len(patch) > 0 {
		newSlug, err = s.store.updateListByID(bson.ObjectIdHex(id), user.ID, newSlug, patch)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, pulpe.ErrListNotFound
			}

			return nil, err
		}
	}

	l, err := s.store.listByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrListNotFound
		}

		return nil, err
	}

	return l.toPulpeList(), nil
}

// ListsByBoard returns all the lists of a given board.
func (s *ListService) ListsByBoard(boardID string) ([]*pulpe.List, error) {
	_, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	ls, err := s.store.listsByBoardID(boardID)
	if err != nil {
		return nil, err
	}

	lists := make([]*pulpe.List, len(ls))
	for i := range ls {
		lists[i] = ls[i].toPulpeList()
	}

	return lists, nil
}

type listStore struct {
	session *Session
}

func (s *listStore) listByOwnerIDAndID(ownerID, id string) (*list, error) {
	var b list

	if !bson.IsObjectIdHex(id) {
		return nil, mgo.ErrNotFound
	}

	query := bson.M{
		"ownerID": ownerID,
		"_id":     bson.ObjectIdHex(id),
	}

	return &b, s.session.db.C(listCol).Find(query).One(&b)
}

func (s *listStore) createList(l *list) error {
	var err error
	col := s.session.db.C(listCol)

	l.Slug, err = resolveSlugAndDo(col, l.OwnerID, "slug", l.Slug, "-", func(slug string) error {
		l.Slug = slug
		return col.Insert(l)
	})

	return err
}

func (s *listStore) deleteListByOwnerIDAndID(ownerID string, id bson.ObjectId) error {
	return s.session.db.C(listCol).Remove(bson.M{
		"_id":     id,
		"ownerID": ownerID,
	})
}

func (s *listStore) deleteListsByBoardID(boardID string) error {
	_, err := s.session.db.C(listCol).RemoveAll(bson.M{
		"boardID": boardID,
	})

	return err
}

func (s *listStore) updateListByID(id bson.ObjectId, ownerID, slug string, patch bson.M) (string, error) {
	col := s.session.db.C(listCol)

	newSlug, err := resolveSlugAndDo(col, ownerID, "slug", slug, "-", func(slug string) error {
		if slug != "" {
			patch["slug"] = slug
		}

		return col.Update(
			bson.M{
				"_id":     id,
				"ownerID": ownerID,
			},
			bson.M{
				"$set":         patch,
				"$currentDate": bson.M{"updatedAt": true},
			})
	})

	return newSlug, err
}

func (s *listStore) listsByBoardID(boardID string) ([]list, error) {
	col := s.session.db.C(listCol)

	var lists []list

	// TODO set a limit
	err := col.Find(bson.M{"boardID": boardID}).Sort("_id").All(&lists)
	return lists, err
}

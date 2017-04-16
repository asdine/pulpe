package mongo

import (
	"time"

	"github.com/Machiel/slugify"
	"github.com/blankrobot/pulpe"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const boardCol = "boards"

// Ensure BoardService implements pulpe.BoardService.
var _ pulpe.BoardService = new(BoardService)

// Board representation stored in MongoDB.
type Board struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	Name      string        `bson:"name"`
	Slug      string        `bson:"slug"`
}

// ToMongoBoard creates a mongo board from a pulpe board.
func ToMongoBoard(p *pulpe.Board) *Board {
	id := bson.NewObjectId()
	p.ID = id.Hex()
	p.CreatedAt = id.Time()

	b := Board{
		ID:        id,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
		Slug:      p.Slug,
	}

	return &b
}

// FromMongoBoard creates a pulpe board from a mongo board.
func FromMongoBoard(b *Board) *pulpe.Board {
	p := pulpe.Board{
		ID:        b.ID.Hex(),
		CreatedAt: b.ID.Time(),
		Name:      b.Name,
		Slug:      b.Slug,
	}

	if b.UpdatedAt != nil {
		t := (*b.UpdatedAt).UTC()
		p.UpdatedAt = &t
	}

	return &p
}

// BoardService represents a service for managing boards.
type BoardService struct {
	session *Session
}

func (s *BoardService) ensureIndexes() error {
	col := s.session.db.C(boardCol)

	// Unique publicID
	index := mgo.Index{
		Key:    []string{"slug"},
		Unique: true,
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateBoard creates a new Board.
func (s *BoardService) CreateBoard(bc *pulpe.BoardCreation) (*pulpe.Board, error) {
	var err error
	col := s.session.db.C(boardCol)

	board := pulpe.Board{
		Name: bc.Name,
		Slug: slugify.Slugify(bc.Name),
	}

	b := ToMongoBoard(&board)

	board.Slug, err = resolveSlugAndDo(col, newBoardRecorder(b), func(rec recorder) error {
		return col.Insert(rec.elem())
	})

	return &board, err
}

// Board returns a Board by slug or ID.
func (s *BoardService) Board(selector string) (*pulpe.Board, error) {
	var b Board

	err := s.session.db.C(boardCol).Find(getSelector(selector)).One(&b)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	return FromMongoBoard(&b), nil
}

// Boards returns all the boards.
func (s *BoardService) Boards() ([]*pulpe.Board, error) {
	var bs []Board

	err := s.session.db.C(boardCol).Find(nil).All(&bs)
	if err != nil {
		return nil, err
	}

	boards := make([]*pulpe.Board, len(bs))
	for i := range bs {
		boards[i] = FromMongoBoard(&bs[i])
	}

	return boards, nil
}

// DeleteBoard deletes a Board by ID, and all of its lists and cards .
func (s *BoardService) DeleteBoard(id string) error {
	if !bson.IsObjectIdHex(id) {
		return pulpe.ErrBoardNotFound
	}

	err := s.session.db.C(boardCol).RemoveId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return pulpe.ErrBoardNotFound
	}

	return err
}

// UpdateBoard updates a Board by ID.
func (s *BoardService) UpdateBoard(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrBoardNotFound
	}

	var b Board
	var err error

	col := s.session.db.C(boardCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		b.Slug = slugify.Slugify(*u.Name)
	}

	if len(patch) == 0 {
		return s.Board(id)
	}

	b.Slug, err = resolveSlugAndDo(col, newBoardRecorder(&b), func(rec recorder) error {
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
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	return s.Board(id)
}

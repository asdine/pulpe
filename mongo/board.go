package mongo

import (
	"encoding/json"
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
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	Name      string        `bson:"name"`
	Slug      string        `bson:"slug"`
	Settings  []byte        `bson:"settings,omitempty"`
}

// ToMongoBoard creates a mongo board from a pulpe board.
func ToMongoBoard(p *pulpe.Board) *Board {
	var id bson.ObjectId

	if p.ID == "" {
		id = bson.NewObjectId()
		p.ID = pulpe.BoardID(id.Hex())
	} else {
		id = bson.ObjectIdHex(string(p.ID))
	}

	b := Board{
		ID:        id,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
		Slug:      p.Slug,
	}

	if p.Settings != nil {
		b.Settings = []byte(*p.Settings)
	}

	return &b
}

// FromMongoBoard creates a pulpe board from a mongo board.
func FromMongoBoard(b *Board) *pulpe.Board {
	p := pulpe.Board{
		ID:        pulpe.BoardID(b.ID.Hex()),
		CreatedAt: b.CreatedAt.UTC(),
		Name:      b.Name,
		Slug:      b.Slug,
		Lists:     []*pulpe.List{},
		Cards:     []*pulpe.Card{},
	}

	if b.UpdatedAt != nil {
		t := (*b.UpdatedAt).UTC()
		p.UpdatedAt = &t
	}

	if len(b.Settings) > 0 {
		s := json.RawMessage(b.Settings)
		p.Settings = &s
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
func (s *BoardService) CreateBoard(b *pulpe.BoardCreate) (*pulpe.Board, error) {
	col := s.session.db.C(boardCol)

	slug := slugify.Slugify(b.Name)

	total, err := col.Find(bson.M{"slug": slug}).Limit(1).Count()
	if err != nil {
		return nil, err
	}

	if total > 0 {
		return nil, pulpe.ErrBoardExists
	}

	board := pulpe.Board{
		CreatedAt: s.session.now,
		Name:      b.Name,
		Slug:      slug,
		Lists:     []*pulpe.List{},
		Cards:     []*pulpe.Card{},
		Settings:  b.Settings,
	}

	return &board, s.session.db.C(boardCol).Insert(ToMongoBoard(&board))
}

// Board returns a Board by ID.
func (s *BoardService) Board(id pulpe.BoardID) (*pulpe.Board, error) {
	var b Board

	err := s.session.db.C(boardCol).FindId(bson.ObjectIdHex(string(id))).One(&b)
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
func (s *BoardService) DeleteBoard(id pulpe.BoardID) error {
	err := s.session.db.C(boardCol).RemoveId(bson.ObjectIdHex(string(id)))
	if err == mgo.ErrNotFound {
		return pulpe.ErrBoardNotFound
	}

	return err
}

// UpdateBoard updates a Board by ID.
func (s *BoardService) UpdateBoard(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
	col := s.session.db.C(boardCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
	}

	if u.Settings != nil {
		patch["settings"] = []byte(*u.Settings)
	}

	if len(patch) == 0 {
		return s.Board(id)
	}

	err := col.UpdateId(
		bson.ObjectIdHex(string(id)),
		bson.M{
			"$set":         patch,
			"$currentDate": bson.M{"updatedAt": true},
		})
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	return s.Board(id)
}

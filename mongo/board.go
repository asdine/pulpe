package mongo

import (
	"encoding/json"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
)

const (
	boardCol = "boards"
)

// Ensure BoardService implements pulpe.BoardService.
var _ pulpe.BoardService = new(BoardService)

// Board representation stored in MongoDB.
type Board struct {
	ID        bson.ObjectId `bson:"_id"`
	PublicID  string        `bson:"publicID"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	Name      string        `bson:"name"`
	Settings  []byte        `bson:"settings,omitempty"`
}

// ToMongoBoard creates a mongo board from a pulpe board.
func ToMongoBoard(p *pulpe.Board) *Board {
	b := Board{
		ID:        bson.NewObjectId(),
		PublicID:  string(p.ID),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Name:      p.Name,
	}

	if p.Settings != nil {
		b.Settings = []byte(*p.Settings)
	}

	return &b
}

// FromMongoBoard creates a pulpe board from a mongo board.
func FromMongoBoard(b *Board) *pulpe.Board {
	p := pulpe.Board{
		ID:        pulpe.BoardID(b.PublicID),
		CreatedAt: b.CreatedAt.UTC(),
		Name:      b.Name,
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

// CreateBoard creates a new Board.
func (s *BoardService) CreateBoard(b *pulpe.BoardCreate) (*pulpe.Board, error) {
	var err error

	board := pulpe.Board{
		Name:      b.Name,
		CreatedAt: s.session.now,
		Lists:     []*pulpe.List{},
		Cards:     []*pulpe.Card{},
		Settings:  b.Settings,
	}

	board.ID, err = pulpe.NewBoardID()
	if err != nil {
		return nil, err
	}

	return &board, s.session.db.C(boardCol).Insert(ToMongoBoard(&board))
}

// Board returns a Board by ID.
func (s *BoardService) Board(id pulpe.BoardID) (*pulpe.Board, error) {
	var b Board

	err := s.session.db.C(boardCol).Find(bson.M{"publicID": string(id)}).One(&b)
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
	return nil, nil
}

// DeleteBoard deletes a Board by ID, and all of its lists and cards .
func (s *BoardService) DeleteBoard(id pulpe.BoardID) error {
	return nil
}

// UpdateBoard updates a Board by ID.
func (s *BoardService) UpdateBoard(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
	return nil, nil
}

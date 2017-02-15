package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/blankrobot/pulpe"
)

const cardCol = "cards"

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// Card representation stored in MongoDB.
type Card struct {
	ID          bson.ObjectId `bson:"_id"`
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   *time.Time    `bson:"updatedAt,omitempty"`
	ListID      bson.ObjectId `bson:"listID"`
	BoardID     bson.ObjectId `bson:"boardID"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Position    float64       `bson:"position"`
}

// ToMongoCard creates a mongo card from a pulpe card.
func ToMongoCard(p *pulpe.Card) *Card {
	var id bson.ObjectId

	if p.ID == "" {
		id = bson.NewObjectId()
		p.ID = id.Hex()
	} else {
		id = bson.ObjectIdHex(p.ID)
	}

	return &Card{
		ID:          id,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		ListID:      bson.ObjectIdHex(p.ListID),
		BoardID:     bson.ObjectIdHex(p.BoardID),
		Name:        p.Name,
		Description: p.Description,
		Position:    p.Position,
	}
}

// FromMongoCard creates a pulpe card from a mongo card.
func FromMongoCard(c *Card) *pulpe.Card {
	p := pulpe.Card{
		ID:          c.ID.Hex(),
		CreatedAt:   c.CreatedAt.UTC(),
		ListID:      c.ListID.Hex(),
		BoardID:     c.BoardID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		Position:    c.Position,
	}

	if c.UpdatedAt != nil {
		t := (*c.UpdatedAt).UTC()
		p.UpdatedAt = &t
	}

	return &p
}

// CardService represents a service for managing cards.
type CardService struct {
	session *Session
}

func (s *CardService) ensureIndexes() error {
	col := s.session.db.C(cardCol)

	// listID
	index := mgo.Index{
		Key:    []string{"listID"},
		Sparse: true,
	}

	err := col.EnsureIndex(index)
	if err != nil {
		return err
	}

	// boardID
	index = mgo.Index{
		Key:    []string{"boardID"},
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateCard creates a new Card.
func (s *CardService) CreateCard(c *pulpe.CardCreate) (*pulpe.Card, error) {
	if c.ListID == "" {
		return nil, pulpe.ErrCardListIDRequired
	}

	if c.BoardID == "" {
		return nil, pulpe.ErrCardBoardIDRequired
	}

	card := pulpe.Card{
		CreatedAt:   s.session.now,
		BoardID:     c.BoardID,
		ListID:      c.ListID,
		Name:        c.Name,
		Description: c.Description,
		Position:    c.Position,
	}

	return &card, s.session.db.C(cardCol).Insert(ToMongoCard(&card))
}

// Card returns a Card by ID.
func (s *CardService) Card(id string) (*pulpe.Card, error) {
	var c Card

	err := s.session.db.C(cardCol).FindId(bson.ObjectIdHex(id)).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrCardNotFound
		}

		return nil, err
	}

	return FromMongoCard(&c), nil
}

// DeleteCard deletes a Card by ID.
func (s *CardService) DeleteCard(id string) error {
	err := s.session.db.C(cardCol).RemoveId(bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return pulpe.ErrCardNotFound
	}

	return err
}

// DeleteCardsByListID deletes all the cards of a list.
func (s *CardService) DeleteCardsByListID(listID string) error {
	_, err := s.session.db.C(cardCol).RemoveAll(bson.M{"listID": bson.ObjectIdHex(listID)})
	return err
}

// DeleteCardsByBoardID deletes all the cards of a board.
func (s *CardService) DeleteCardsByBoardID(boardID string) error {
	_, err := s.session.db.C(cardCol).RemoveAll(bson.M{"boardID": bson.ObjectIdHex(boardID)})
	return err
}

// UpdateCard updates a Card by ID.
func (s *CardService) UpdateCard(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
	col := s.session.db.C(cardCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
	}

	if u.Description != nil {
		patch["description"] = *u.Description
	}

	if u.Position != nil {
		patch["position"] = *u.Position
	}

	if len(patch) == 0 {
		return s.Card(id)
	}

	err := col.UpdateId(
		bson.ObjectIdHex(id),
		bson.M{
			"$set":         patch,
			"$currentDate": bson.M{"updatedAt": true},
		})
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrCardNotFound
		}

		return nil, err
	}

	return s.Card(id)
}

// CardsByBoard returns Cards by board ID.
func (s *CardService) CardsByBoard(boardID string) ([]*pulpe.Card, error) {
	var cards []Card

	// TODO set a limit
	err := s.session.db.C(cardCol).Find(bson.M{"boardID": bson.ObjectIdHex(boardID)}).All(&cards)
	if err != nil {
		return nil, err
	}

	list := make([]*pulpe.Card, len(cards))
	for i := range cards {
		list[i] = FromMongoCard(&cards[i])
	}

	return list, nil
}

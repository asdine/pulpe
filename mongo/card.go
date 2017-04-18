package mongo

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Machiel/slugify"
	"github.com/blankrobot/pulpe"
)

const cardCol = "cards"

const codeConflict = 11000

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// Card representation stored in MongoDB.
type Card struct {
	ID          bson.ObjectId `bson:"_id"`
	UpdatedAt   *time.Time    `bson:"updatedAt,omitempty"`
	ListID      bson.ObjectId `bson:"listID"`
	BoardID     bson.ObjectId `bson:"boardID"`
	Name        string        `bson:"name"`
	Slug        string        `bson:"slug"`
	Description string        `bson:"description"`
	Position    float64       `bson:"position"`
}

// ToMongoCard creates a mongo card from a pulpe card.
func ToMongoCard(p *pulpe.Card) *Card {
	id := bson.NewObjectId()
	p.ID = id.Hex()
	p.CreatedAt = id.Time().UTC()

	return &Card{
		ID:          id,
		UpdatedAt:   p.UpdatedAt,
		ListID:      bson.ObjectIdHex(p.ListID),
		BoardID:     bson.ObjectIdHex(p.BoardID),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Position:    p.Position,
	}
}

// FromMongoCard creates a pulpe card from a mongo card.
func FromMongoCard(c *Card) *pulpe.Card {
	p := pulpe.Card{
		ID:          c.ID.Hex(),
		CreatedAt:   c.ID.Time().UTC(),
		ListID:      c.ListID.Hex(),
		BoardID:     c.BoardID.Hex(),
		Name:        c.Name,
		Slug:        c.Slug,
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

	// boardID and slug
	index := mgo.Index{
		Key:    []string{"boardID", "slug"},
		Unique: true,
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateCard creates a new Card.
func (s *CardService) CreateCard(cc *pulpe.CardCreation) (*pulpe.Card, error) {
	var err error

	// generate slug
	slug := slugify.Slugify(cc.Name)

	// create mongo card and pulpe card
	card := pulpe.Card{
		BoardID:     cc.BoardID,
		ListID:      cc.ListID,
		Name:        cc.Name,
		Slug:        slug,
		Description: cc.Description,
		Position:    cc.Position,
	}

	c := ToMongoCard(&card)
	col := s.session.db.C(cardCol)

	card.Slug, err = resolveSlugAndDo(col, "slug", c.Slug, "-", func(slug string) error {
		c.Slug = slug
		return col.Insert(c)
	})

	return &card, err
}

// Card returns a Card by ID.
func (s *CardService) Card(id string) (*pulpe.Card, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrCardNotFound
	}

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
	if !bson.IsObjectIdHex(id) {
		return pulpe.ErrCardNotFound
	}

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
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrCardNotFound
	}

	var err error
	var c Card

	col := s.session.db.C(cardCol)

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		c.Slug = slugify.Slugify(*u.Name)
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

	c.Slug, err = resolveSlugAndDo(col, "slug", c.Slug, "-", func(slug string) error {
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
	err := s.session.db.C(cardCol).Find(bson.M{"boardID": bson.ObjectIdHex(boardID)}).Sort("_id").All(&cards)
	if err != nil {
		return nil, err
	}

	list := make([]*pulpe.Card, len(cards))
	for i := range cards {
		list[i] = FromMongoCard(&cards[i])
	}

	return list, nil
}

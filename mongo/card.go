package mongo

import (
	"time"

	"github.com/Machiel/slugify"
	"github.com/blankrobot/pulpe"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const cardCol = "cards"

// Ensure CardService implements pulpe.CardService.
var _ pulpe.CardService = new(CardService)

// card representation stored in MongoDB.
type card struct {
	ID          bson.ObjectId `bson:"_id"`
	UpdatedAt   *time.Time    `bson:"updatedAt,omitempty"`
	OwnerID     string        `bson:"ownerID"`
	ListID      string        `bson:"listID"`
	BoardID     string        `bson:"boardID"`
	Name        string        `bson:"name"`
	Slug        string        `bson:"slug"`
	Description string        `bson:"description"`
	Position    float64       `bson:"position"`
}

// toPulpeCard creates a pulpe card from a mongo card.
func (c *card) toPulpeCard() *pulpe.Card {
	p := pulpe.Card{
		ID:          c.ID.Hex(),
		CreatedAt:   c.ID.Time().UTC(),
		OwnerID:     c.OwnerID,
		ListID:      c.ListID,
		BoardID:     c.BoardID,
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
	store   cardStore
}

func (s *CardService) ensureIndexes() error {
	col := s.session.db.C(cardCol)

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

// CreateCard creates a new Card.
func (s *CardService) CreateCard(listID string, cc *pulpe.CardCreation) (*pulpe.Card, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	list, err := s.session.ListService().List(listID)
	if err != nil {
		return nil, err
	}

	if list.OwnerID != user.ID {
		return nil, pulpe.ErrListNotFound
	}

	c := card{
		ID:          bson.NewObjectId(),
		OwnerID:     user.ID,
		BoardID:     list.BoardID,
		ListID:      listID,
		Name:        cc.Name,
		Slug:        slugify.Slugify(cc.Name),
		Description: cc.Description,
		Position:    cc.Position,
	}

	err = s.store.createCard(&c)
	if err != nil {
		return nil, err
	}

	return c.toPulpeCard(), err
}

// Card returns a Card.
func (s *CardService) Card(id string) (*pulpe.Card, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrCardNotFound
	}

	c, err := s.store.cardByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrCardNotFound
		}

		return nil, err
	}

	return c.toPulpeCard(), nil
}

// DeleteCard deletes a Card by ID.
func (s *CardService) DeleteCard(id string) error {
	user, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	if !bson.IsObjectIdHex(id) {
		return pulpe.ErrCardNotFound
	}

	err = s.store.deleteCardByID(user.ID, bson.ObjectIdHex(id))
	if err == mgo.ErrNotFound {
		return pulpe.ErrCardNotFound
	}

	return err
}

// UpdateCard updates a Card by ID.
func (s *CardService) UpdateCard(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrCardNotFound
	}

	var newSlug string

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		newSlug = slugify.Slugify(*u.Name)
	}

	if u.Description != nil {
		patch["description"] = *u.Description
	}

	if u.Position != nil {
		patch["position"] = *u.Position
	}

	if u.ListID != nil {
		list, err := s.session.ListService().List(*u.ListID)
		if err != nil {
			return nil, err
		}

		if list.OwnerID != user.ID {
			return nil, pulpe.ErrListNotFound
		}
		patch["listID"] = *u.ListID
	}

	if len(patch) > 0 {
		newSlug, err = s.store.updateCardByID(bson.ObjectIdHex(id), user.ID, newSlug, patch)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, pulpe.ErrCardNotFound
			}

			return nil, err
		}
	}

	c, err := s.store.cardByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrCardNotFound
		}

		return nil, err
	}

	return c.toPulpeCard(), nil
}

// CardsByBoard returns Cards by board ID.
func (s *CardService) CardsByBoard(boardID string) ([]*pulpe.Card, error) {
	_, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	cs, err := s.store.cardsByBoardID(boardID)
	if err != nil {
		return nil, err
	}

	cards := make([]*pulpe.Card, len(cs))
	for i := range cs {
		cards[i] = cs[i].toPulpeCard()
	}

	return cards, nil
}

// DeleteCardsByListID deletes all the cards of a list.
func (s *CardService) DeleteCardsByListID(listID string) error {
	user, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	return s.store.deleteCardsByListID(user.ID, listID)
}

// DeleteCardsByBoardID deletes all the cards of a board.
func (s *CardService) DeleteCardsByBoardID(boardID string) error {
	user, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	return s.store.deleteCardsByBoardID(user.ID, boardID)
}

type cardStore struct {
	session *Session
}

func (s *cardStore) createCard(c *card) error {
	var err error
	col := s.session.db.C(cardCol)

	c.Slug, err = resolveSlugAndDo(col, c.OwnerID, "slug", c.Slug, "-", func(slug string) error {
		c.Slug = slug
		return col.Insert(c)
	})

	return err
}

func (s *cardStore) cardByOwnerIDAndID(ownerID, id string) (*card, error) {
	var c card

	if !bson.IsObjectIdHex(id) {
		return nil, mgo.ErrNotFound
	}

	query := bson.M{
		"ownerID": ownerID,
		"_id":     bson.ObjectIdHex(id),
	}

	return &c, s.session.db.C(cardCol).Find(query).One(&c)
}

func (s *cardStore) deleteCardByID(ownerID string, id bson.ObjectId) error {
	return s.session.db.C(cardCol).Remove(bson.M{
		"_id":     id,
		"ownerID": ownerID,
	})
}

func (s *cardStore) deleteCardsByBoardID(ownerID, boardID string) error {
	_, err := s.session.db.C(cardCol).RemoveAll(bson.M{
		"ownerID": ownerID,
		"boardID": boardID,
	})

	return err
}

func (s *cardStore) deleteCardsByListID(ownerID, listID string) error {
	_, err := s.session.db.C(cardCol).RemoveAll(bson.M{
		"ownerID": ownerID,
		"listID":  listID,
	})

	return err
}

func (s *cardStore) updateCardByID(id bson.ObjectId, ownerID, slug string, patch bson.M) (string, error) {
	col := s.session.db.C(cardCol)

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

func (s *cardStore) cardsByBoardID(boardID string) ([]card, error) {
	col := s.session.db.C(cardCol)

	var cards []card

	// TODO set a limit
	err := col.Find(bson.M{"boardID": boardID}).Sort("_id").All(&cards)
	return cards, err
}

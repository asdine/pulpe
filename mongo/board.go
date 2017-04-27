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

// board representation stored in MongoDB.
type board struct {
	ID        bson.ObjectId `bson:"_id"`
	UpdatedAt *time.Time    `bson:"updatedAt,omitempty"`
	Name      string        `bson:"name"`
	Slug      string        `bson:"slug"`
	OwnerID   string        `bson:"ownerID"`
}

// toPulpeBoard creates a pulpe board from a mongo board.
func (b *board) toPulpeBoard() *pulpe.Board {
	p := pulpe.Board{
		ID:        b.ID.Hex(),
		CreatedAt: b.ID.Time().UTC(),
		Name:      b.Name,
		Slug:      b.Slug,
		OwnerID:   b.OwnerID,
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
	store   boardStore
}

func (s *BoardService) ensureIndexes() error {
	col := s.session.db.C(boardCol)

	// Unique slug
	index := mgo.Index{
		Key:    []string{"ownerID", "slug"},
		Unique: true,
		Sparse: true,
	}

	return col.EnsureIndex(index)
}

// CreateBoard creates a new Board.
func (s *BoardService) CreateBoard(bc *pulpe.BoardCreation) (*pulpe.Board, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	b := board{
		ID:      bson.NewObjectId(),
		Name:    bc.Name,
		Slug:    slugify.Slugify(bc.Name),
		OwnerID: user.ID,
	}

	err = s.store.createBoard(&b)
	if err != nil {
		return nil, err
	}

	return b.toPulpeBoard(), nil
}

// Board returns a Board by id.
func (s *BoardService) Board(id string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	b, err := s.store.boardByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	board := b.toPulpeBoard()

	var opts pulpe.BoardGetOptions

	for i := range options {
		options[i](&opts)
	}

	if opts.WithLists {
		board.Lists, err = s.session.ListService().ListsByBoard(board.ID)
		if err != nil {
			return nil, err
		}
	}

	if opts.WithCards {
		board.Cards, err = s.session.CardService().CardsByBoard(board.ID)
		if err != nil {
			return nil, err
		}
	}

	return board, nil
}

// BoardByOwnerAndSlug returns a Board by owner and slug.
func (s *BoardService) BoardByOwnerAndSlug(owner, slug string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	if user.Login != owner {
		return nil, pulpe.ErrBoardNotFound
	}

	b, err := s.store.boardByOwnerIDAndSlug(user.ID, slug)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	if b.OwnerID != user.ID {
		return nil, pulpe.ErrBoardNotFound
	}

	board := b.toPulpeBoard()

	var opts pulpe.BoardGetOptions

	for i := range options {
		options[i](&opts)
	}

	if opts.WithLists {
		board.Lists, err = s.session.ListService().ListsByBoard(board.ID)
		if err != nil {
			return nil, err
		}
	}

	if opts.WithCards {
		board.Cards, err = s.session.CardService().CardsByBoard(board.ID)
		if err != nil {
			return nil, err
		}
	}

	return board, nil
}

// Boards returns all the boards of the authenticated user.
func (s *BoardService) Boards() ([]*pulpe.Board, error) {
	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	bs, err := s.store.boardsByOwnerID(user.ID)
	if err != nil {
		return nil, err
	}

	boards := make([]*pulpe.Board, len(bs))
	for i := range bs {
		boards[i] = bs[i].toPulpeBoard()
	}

	return boards, nil
}

// DeleteBoard deletes a Board and its related cards and lists.
func (s *BoardService) DeleteBoard(id string) error {
	user, err := s.session.Authenticate()
	if err != nil {
		return err
	}

	b, err := s.store.boardByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return pulpe.ErrBoardNotFound
		}

		return err
	}

	err = s.store.deleteBoardByID(b.ID)
	if err == mgo.ErrNotFound {
		return pulpe.ErrBoardNotFound
	}

	err = s.session.ListService().DeleteListsByBoardID(b.ID.Hex())
	if err != nil {
		return err
	}

	return s.session.CardService().DeleteCardsByBoardID(b.ID.Hex())
}

// UpdateBoard updates a Board.
func (s *BoardService) UpdateBoard(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, pulpe.ErrBoardNotFound
	}

	user, err := s.session.Authenticate()
	if err != nil {
		return nil, err
	}

	var newSlug string

	patch := make(bson.M)
	if u.Name != nil {
		patch["name"] = *u.Name
		newSlug = slugify.Slugify(*u.Name)
	}

	if len(patch) > 0 {
		newSlug, err = s.store.updateBoardByID(bson.ObjectIdHex(id), user.ID, newSlug, patch)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, pulpe.ErrBoardNotFound
			}

			return nil, err
		}
	}

	b, err := s.store.boardByOwnerIDAndID(user.ID, id)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, pulpe.ErrBoardNotFound
		}

		return nil, err
	}

	return b.toPulpeBoard(), nil
}

type boardStore struct {
	session *Session
}

func (s *boardStore) boardByOwnerIDAndID(ownerID, id string) (*board, error) {
	var b board

	if !bson.IsObjectIdHex(id) {
		return nil, mgo.ErrNotFound
	}

	query := bson.M{
		"ownerID": ownerID,
		"_id":     bson.ObjectIdHex(id),
	}

	return &b, s.session.db.C(boardCol).Find(query).One(&b)
}

func (s *boardStore) boardByOwnerIDAndSlug(ownerID, slug string) (*board, error) {
	var b board

	query := bson.M{
		"ownerID": ownerID,
		"slug":    slug,
	}

	return &b, s.session.db.C(boardCol).Find(query).One(&b)
}

func (s *boardStore) createBoard(b *board) error {
	var err error
	col := s.session.db.C(boardCol)

	b.Slug, err = resolveSlugAndDo(col, b.OwnerID, "slug", b.Slug, "-", func(slug string) error {
		b.Slug = slug
		return col.Insert(b)
	})

	return err
}

func (s *boardStore) boardsByOwnerID(ownerID string) ([]board, error) {
	var bs []board

	return bs, s.session.db.C(boardCol).Find(bson.M{"ownerID": ownerID}).Sort("_id").All(&bs)
}

func (s *boardStore) deleteBoardByID(id bson.ObjectId) error {
	return s.session.db.C(boardCol).RemoveId(id)
}

func (s *boardStore) updateBoardByID(id bson.ObjectId, ownerID, slug string, patch bson.M) (string, error) {
	col := s.session.db.C(boardCol)

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

package pulpe_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestCardCreate_Validate(t *testing.T) {
	c := mock.NewClient()
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		if id != "XXX" {
			return nil, pulpe.ErrBoardNotFound
		}
		return new(pulpe.Board), nil
	}

	c.ListService.ListFn = func(id pulpe.ListID) (*pulpe.List, error) {
		if id != "XXX" {
			return nil, pulpe.ErrListNotFound
		}
		return new(pulpe.List), nil
	}

	t.Run("Empty", func(t *testing.T) {
		var cc pulpe.CardCreate
		err := cc.Validate(c.Connect())
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name: "Card name",
		}
		require.Error(t, cc.Validate(c.Connect()))
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name:    "      ",
			BoardID: "    ",
			ListID:  "    ",
		}
		require.Error(t, cc.Validate(c.Connect()))
	})

	t.Run("Valid", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name:    "Card name",
			BoardID: "XXX",
			ListID:  "XXX",
		}
		require.NoError(t, cc.Validate(c.Connect()))
	})

	t.Run("UnknownBoard", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name:    "Card name",
			BoardID: "ZZZ",
			ListID:  "XXX",
		}
		require.Error(t, cc.Validate(c.Connect()))
	})

	t.Run("UnknownList", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name:    "Card name",
			BoardID: "XXX",
			ListID:  "ZZZ",
		}
		require.Error(t, cc.Validate(c.Connect()))
	})

	t.Run("NegativePosition", func(t *testing.T) {
		cc := pulpe.CardCreate{
			Name:     "Card name",
			BoardID:  "XXX",
			ListID:   "XXX",
			Position: -10.0,
		}
		require.Error(t, cc.Validate(c.Connect()))
	})
}

func TestCardUpdate_Validate(t *testing.T) {
	name := "Card name"
	emptyName := ""
	spaces := "    "
	zeroFloat := 0.0
	negativeFloat := -10.0
	positiveFloat := 10.0

	t.Run("Empty", func(t *testing.T) {
		var cc pulpe.CardUpdate
		err := cc.Validate()
		require.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		cc := pulpe.CardUpdate{
			Name:        &name,
			Description: &name,
			Position:    &positiveFloat,
		}
		require.NoError(t, cc.Validate())
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		cc := pulpe.CardUpdate{
			Name:        &spaces,
			Description: &spaces,
		}
		require.Error(t, cc.Validate())
	})

	t.Run("ZeroValues", func(t *testing.T) {
		cc := pulpe.CardUpdate{
			Name:        &emptyName,
			Description: &emptyName,
			Position:    &zeroFloat,
		}
		require.Error(t, cc.Validate())
	})

	t.Run("NegativePosition", func(t *testing.T) {
		cc := pulpe.CardUpdate{
			Name:        &name,
			Description: &name,
			Position:    &negativeFloat,
		}
		require.Error(t, cc.Validate())
	})
}

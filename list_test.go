package pulpe_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestListCreate_Validate(t *testing.T) {
	c := mock.NewClient()
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		if id != "XXX" {
			return nil, pulpe.ErrBoardNotFound
		}
		return new(pulpe.Board), nil
	}

	t.Run("Empty", func(t *testing.T) {
		var l pulpe.ListCreate
		err := l.Validate(c.Connect())
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name: "list name",
		}
		require.Error(t, l.Validate(c.Connect()))
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name:    "      ",
			BoardID: "    ",
		}
		require.Error(t, l.Validate(c.Connect()))
	})

	t.Run("Valid", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name:    "list name",
			BoardID: "XXX",
		}
		require.NoError(t, l.Validate(c.Connect()))
	})

	t.Run("UnknownBoard", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name:    "list name",
			BoardID: "ZZZ",
		}
		require.Error(t, l.Validate(c.Connect()))
	})
}

func TestListUpdate_Validate(t *testing.T) {
	name := "list name"
	emptyName := ""
	spaces := "    "

	t.Run("Empty", func(t *testing.T) {
		var l pulpe.ListUpdate
		err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("ValidName", func(t *testing.T) {
		l := pulpe.ListUpdate{
			Name: &name,
		}
		require.NoError(t, l.Validate())
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := pulpe.ListUpdate{
			Name: &spaces,
		}
		require.Error(t, l.Validate())
	})

	t.Run("EmptyName", func(t *testing.T) {
		l := pulpe.ListUpdate{
			Name: &emptyName,
		}
		require.Error(t, l.Validate())
	})
}

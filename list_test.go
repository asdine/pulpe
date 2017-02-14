package pulpe_test

import (
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func TestListCreate_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var l pulpe.ListCreate
		err := l.Validate()
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name: "name",
		}
		require.Error(t, l.Validate())
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name:    "      ",
			BoardID: "    ",
		}
		require.Error(t, l.Validate())
	})

	t.Run("Valid", func(t *testing.T) {
		l := pulpe.ListCreate{
			Name:    "name",
			BoardID: "boardID",
		}
		require.NoError(t, l.Validate())
	})
}

func TestListUpdate_Validate(t *testing.T) {
	name := "name"
	emptyName := ""
	spaces := "    "

	t.Run("Empty", func(t *testing.T) {
		var l pulpe.ListUpdate
		err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("Valid name", func(t *testing.T) {
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

package pulpe_test

import (
	"encoding/json"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func TestBoardCreate_Validate(t *testing.T) {
	raw := json.RawMessage("{}")

	t.Run("Empty", func(t *testing.T) {
		var b pulpe.BoardCreate
		err := b.Validate()
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "    name   ",
		}
		require.NoError(t, b.Validate())
		require.Equal(t, "name", b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "      ",
		}
		require.Error(t, b.Validate())
	})

	t.Run("WithSettings", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name:     "name",
			Settings: &raw,
		}
		require.NoError(t, b.Validate())
	})
}

func TestBoardUpdate_Validate(t *testing.T) {
	emptyName := ""
	spaces := "    "
	raw := json.RawMessage("{}")

	t.Run("Empty", func(t *testing.T) {
		var b pulpe.BoardUpdate
		err := b.Validate()
		require.NoError(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		name := "   name   "
		b := pulpe.BoardUpdate{
			Name: &name,
		}
		require.NoError(t, b.Validate())
		require.Equal(t, "name", *b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := pulpe.BoardUpdate{
			Name: &spaces,
		}
		require.Error(t, b.Validate())
	})

	t.Run("WithSettings", func(t *testing.T) {
		name := "   name   "
		b := pulpe.BoardUpdate{
			Name:     &name,
			Settings: &raw,
		}
		require.NoError(t, b.Validate())
	})

	t.Run("EmptyName", func(t *testing.T) {
		b := pulpe.BoardUpdate{
			Name: &emptyName,
		}
		require.Error(t, b.Validate())
	})
}

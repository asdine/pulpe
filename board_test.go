package pulpe_test

import (
	"encoding/json"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/stretchr/testify/require"
)

func TestBoardCreate_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var b pulpe.BoardCreate
		err := b.Validate()
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		b := pulpe.BoardCreate{
			Name: "name",
		}
		require.NoError(t, b.Validate())

		raw := json.RawMessage("{}")
		b.Settings = &raw
		require.NoError(t, b.Validate())
	})

	t.Run("WithSettings", func(t *testing.T) {
		raw := json.RawMessage("{}")
		b := pulpe.BoardCreate{
			Name:     "name",
			Settings: &raw,
		}
		require.NoError(t, b.Validate())
	})
}

package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blankrobot/pulpe"
	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestBoardHandler_Boards(t *testing.T) {
	t.Run("OK", testBoardHandler_Boards_OK)
	t.Run("Internal error", testBoardHandler_Boards_InternalError)
}

func testBoardHandler_Boards_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardsFn = func() ([]*pulpe.Board, error) {
		s := json.RawMessage([]byte(`{"a": "b"}`))
		return []*pulpe.Board{
			&pulpe.Board{ID: "id", Name: "name", Slug: "slug", CreatedAt: mock.Now, UpdatedAt: &mock.Now, Lists: []*pulpe.List{}, Cards: []*pulpe.Card{}, Settings: &s},
		}, nil
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `[
    {
			"id": "id",
			"name": "name",
			"slug": "slug",
      "createdAt": `+string(date)+`,
      "updatedAt": `+string(date)+`,
			"lists": [],
			"cards": [],
      "settings": {
        "a": "b"
      }
	  }
  ]`, w.Body.String())
}

func testBoardHandler_Boards_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardsFn = func() ([]*pulpe.Board, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardsInvoked)
}

func TestBoardHandler_CreateBoard(t *testing.T) {
	t.Run("OK", testBoardHandler_CreateBoard_OK)
	t.Run("OKNoSettings", testBoardHandler_CreateBoard_OK_NoSettings)
	t.Run("ErrInvalidJSON", testBoardHandler_CreateBoard_ErrInvalidJSON)
	t.Run("ValidationError", testBoardHandler_CreateBoard_ValidationError)
	t.Run("ErrInternal", testBoardHandler_CreateBoard_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
	t.Run("Conflict", testBoardHandler_CreateBoard_WithResponse(t, http.StatusConflict, pulpe.ErrBoardExists))
}

func testBoardHandler_CreateBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.CreateBoardFn = func(c *pulpe.BoardCreate) (*pulpe.Board, error) {
		require.Equal(t, "name", c.Name)
		require.JSONEq(t, `{"a": "b"}`, string(*c.Settings))

		return &pulpe.Board{
			ID:        "123",
			CreatedAt: mock.Now,
			Name:      c.Name,
			Lists:     []*pulpe.List{},
			Cards:     []*pulpe.Card{},
			Settings:  c.Settings,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{
    "name": "name",
    "settings": {
      "a": "b"
    }
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.CreateBoardInvoked)
}

func testBoardHandler_CreateBoard_OK_NoSettings(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.CreateBoardFn = func(c *pulpe.BoardCreate) (*pulpe.Board, error) {
		require.Equal(t, "name", c.Name)
		require.JSONEq(t, `{}`, string(*c.Settings))

		return &pulpe.Board{
			ID:        "123",
			CreatedAt: mock.Now,
			Name:      c.Name,
			Lists:     []*pulpe.List{},
			Cards:     []*pulpe.Card{},
			Settings:  c.Settings,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{
    "name": "name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.CreateBoardInvoked)
}

func testBoardHandler_CreateBoard_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid json"}`, w.Body.String())
}

func testBoardHandler_CreateBoard_ValidationError(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "validation error", "fields": {"name": "cannot be blank"}}`, w.Body.String())
}

func testBoardHandler_CreateBoard_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		// Mock service.
		c.BoardService.CreateBoardFn = func(Board *pulpe.BoardCreate) (*pulpe.Board, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{"name": "name"}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
		require.True(t, c.BoardService.CreateBoardInvoked)
	}
}

func TestBoardHandler_Board(t *testing.T) {
	t.Run("OK", testBoardHandler_Board_OK)
	t.Run("Not found", testBoardHandler_Board_NotFound)
	t.Run("Internal error", testBoardHandler_Board_InternalError)
	t.Run("List Internal error", testBoardHandler_Board_ListInternalError)
	t.Run("Card Internal error", testBoardHandler_Board_CardInternalError)
}

func testBoardHandler_Board_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return new(pulpe.Board), nil
	}

	c.ListService.ListsByBoardFn = func(id pulpe.BoardID) ([]*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		return nil, nil
	}

	c.CardService.CardsByBoardFn = func(id pulpe.BoardID) ([]*pulpe.Card, error) {
		require.Equal(t, "XXX", string(id))
		return nil, nil
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.BoardInvoked)
	require.True(t, c.ListService.ListsByBoardInvoked)
	require.True(t, c.CardService.CardsByBoardInvoked)
}

func testBoardHandler_Board_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{}`, w.Body.String())
	require.True(t, c.BoardService.BoardInvoked)
}

func testBoardHandler_Board_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardInvoked)
}

func testBoardHandler_Board_ListInternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return &pulpe.Board{ID: id, Name: "name", CreatedAt: mock.Now, UpdatedAt: &mock.Now}, nil
	}

	c.ListService.ListsByBoardFn = func(id pulpe.BoardID) ([]*pulpe.List, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardInvoked)
	require.True(t, c.ListService.ListsByBoardInvoked)
}

func testBoardHandler_Board_CardInternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id pulpe.BoardID) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return new(pulpe.Board), nil
	}

	c.ListService.ListsByBoardFn = func(id pulpe.BoardID) ([]*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		return []*pulpe.List{}, nil
	}

	c.CardService.CardsByBoardFn = func(id pulpe.BoardID) ([]*pulpe.Card, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardInvoked)
	require.True(t, c.ListService.ListsByBoardInvoked)
	require.True(t, c.CardService.CardsByBoardInvoked)
}

func TestBoardHandler_DeleteBoard(t *testing.T) {
	t.Run("OK", testBoardHandler_DeleteBoard_OK)
	t.Run("Not found", testBoardHandler_DeleteBoard_NotFound)
	t.Run("Internal error on delete board", testBoardHandler_DeleteBoard_InternalErrorOnDeleteBoard)
	t.Run("Internal error on delete lists by board id", testBoardHandler_DeleteBoard_InternalErrorOnDeleteListsByBoardID)
	t.Run("Internal error on delete cards by board id", testBoardHandler_DeleteBoard_InternalErrorOnDeleteCardsByBoardID)
}

func testBoardHandler_DeleteBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	byBoardID := func(id pulpe.BoardID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.BoardService.DeleteBoardFn = byBoardID
	c.ListService.DeleteListsByBoardIDFn = byBoardID
	c.CardService.DeleteCardsByBoardIDFn = byBoardID

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.True(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.True(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func testBoardHandler_DeleteBoard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.DeleteBoardFn = func(id pulpe.BoardID) error {
		return pulpe.ErrBoardNotFound
	}

	byBoardID := func(id pulpe.BoardID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}
	c.ListService.DeleteListsByBoardIDFn = byBoardID
	c.CardService.DeleteCardsByBoardIDFn = byBoardID

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{}`, w.Body.String())
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.False(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.False(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func testBoardHandler_DeleteBoard_InternalErrorOnDeleteBoard(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.DeleteBoardFn = func(id pulpe.BoardID) error {
		return errors.New("unexpected error")
	}

	byBoardID := func(id pulpe.BoardID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}
	c.ListService.DeleteListsByBoardIDFn = byBoardID
	c.CardService.DeleteCardsByBoardIDFn = byBoardID

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.False(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.False(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func testBoardHandler_DeleteBoard_InternalErrorOnDeleteListsByBoardID(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	byBoardID := func(id pulpe.BoardID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.BoardService.DeleteBoardFn = byBoardID
	c.ListService.DeleteListsByBoardIDFn = func(id pulpe.BoardID) error {
		return errors.New("unexpected error")
	}
	c.CardService.DeleteCardsByBoardIDFn = byBoardID

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.True(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.False(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func testBoardHandler_DeleteBoard_InternalErrorOnDeleteCardsByBoardID(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	byBoardID := func(id pulpe.BoardID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.BoardService.DeleteBoardFn = byBoardID
	c.ListService.DeleteListsByBoardIDFn = byBoardID
	c.CardService.DeleteCardsByBoardIDFn = func(id pulpe.BoardID) error {
		return errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.True(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.True(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func TestBoardHandler_UpdateBoard(t *testing.T) {
	t.Run("OK", testBoardHandler_UpdateBoard_OK)
	t.Run("ErrInvalidJSON", testBoardHandler_UpdateBoard_ErrInvalidJSON)
	t.Run("Not found", testBoardHandler_UpdateBoard_NotFound)
	t.Run("Internal error", testBoardHandler_UpdateBoard_InternalError)
}

func testBoardHandler_UpdateBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.UpdateBoardFn = func(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)
		require.JSONEq(t, `{"a": "b"}`, string(*u.Settings))

		return new(pulpe.Board), nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name",
    "settings": {
      "a": "b"
    }
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid json"}`, w.Body.String())
}

func testBoardHandler_UpdateBoard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.UpdateBoardFn = func(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{}`, w.Body.String())
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.UpdateBoardFn = func(id pulpe.BoardID, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

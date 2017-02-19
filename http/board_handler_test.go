package http_test

import (
	"bytes"
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
		return []*pulpe.Board{
			&pulpe.Board{ID: "id", Name: "name", Slug: "slug", CreatedAt: mock.Now, UpdatedAt: &mock.Now},
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
      "updatedAt": `+string(date)+`
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
	t.Run("ErrInvalidJSON", testBoardHandler_CreateBoard_ErrInvalidJSON)
	t.Run("ValidationError", testBoardHandler_CreateBoard_ValidationError)
	t.Run("ErrInternal", testBoardHandler_CreateBoard_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testBoardHandler_CreateBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.CreateBoardFn = func(c *pulpe.BoardCreate) (*pulpe.Board, error) {
		require.Equal(t, "name", c.Name)

		return &pulpe.Board{
			ID:        "123",
			CreatedAt: mock.Now,
			Name:      c.Name,
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
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testBoardHandler_CreateBoard_ValidationError(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "validation error", "fields": {"name": ["non zero value required"]}}`, w.Body.String())
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
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return &pulpe.Board{
			ID: "XXX",
		}, nil
	}

	c.ListService.ListsByBoardFn = func(id string) ([]*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		return nil, nil
	}

	c.CardService.CardsByBoardFn = func(id string) ([]*pulpe.Card, error) {
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
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.BoardService.BoardInvoked)
}

func testBoardHandler_Board_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
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
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return &pulpe.Board{ID: id, Name: "name", CreatedAt: mock.Now, UpdatedAt: &mock.Now}, nil
	}

	c.ListService.ListsByBoardFn = func(id string) ([]*pulpe.List, error) {
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
	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		return &pulpe.Board{
			ID: "XXX",
		}, nil
	}

	c.ListService.ListsByBoardFn = func(id string) ([]*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		return []*pulpe.List{}, nil
	}

	c.CardService.CardsByBoardFn = func(id string) ([]*pulpe.Card, error) {
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
	byBoardID := func(id string) error {
		require.Equal(t, "XXX", id)
		return nil
	}

	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", id)
		return &pulpe.Board{ID: "XXX"}, nil
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
	c.BoardService.DeleteBoardFn = func(id string) error {
		return pulpe.ErrBoardNotFound
	}

	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", id)
		return &pulpe.Board{ID: "XXX"}, nil
	}

	byBoardID := func(id string) error {
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
	require.True(t, c.BoardService.DeleteBoardInvoked)
	require.False(t, c.ListService.DeleteListsByBoardIDInvoked)
	require.False(t, c.CardService.DeleteCardsByBoardIDInvoked)
}

func testBoardHandler_DeleteBoard_InternalErrorOnDeleteBoard(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.DeleteBoardFn = func(id string) error {
		return errors.New("unexpected error")
	}

	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", id)
		return &pulpe.Board{ID: "XXX"}, nil
	}

	byBoardID := func(id string) error {
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
	byBoardID := func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", id)
		return &pulpe.Board{ID: "XXX"}, nil
	}

	c.BoardService.DeleteBoardFn = byBoardID
	c.ListService.DeleteListsByBoardIDFn = func(id string) error {
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
	byBoardID := func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
		require.Equal(t, "XXX", id)
		return &pulpe.Board{ID: "XXX"}, nil
	}

	c.BoardService.DeleteBoardFn = byBoardID
	c.ListService.DeleteListsByBoardIDFn = byBoardID
	c.CardService.DeleteCardsByBoardIDFn = func(id string) error {
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
	t.Run("Validation error", testBoardHandler_UpdateBoard_ValidationError)
	t.Run("Internal error", testBoardHandler_UpdateBoard_InternalError)
}

func testBoardHandler_UpdateBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)

		return new(pulpe.Board), nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
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
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testBoardHandler_UpdateBoard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_ValidationError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/boards/XXX", bytes.NewReader([]byte(`{
    "name": "       "
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
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

func TestBoardCreateRequest_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var b pulpeHttp.BoardCreateRequest
		_, err := b.Validate()
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		b := pulpeHttp.BoardCreateRequest{
			Name: "    board name   ",
		}
		_, err := b.Validate()
		require.NoError(t, err)
		require.Equal(t, "board name", b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := pulpeHttp.BoardCreateRequest{
			Name: "      ",
		}
		_, err := b.Validate()
		require.Error(t, err)
	})
}

func TestBoardUpdateRequest_Validate(t *testing.T) {
	emptyName := ""
	spaces := "    "

	t.Run("Empty", func(t *testing.T) {
		var b pulpeHttp.BoardUpdateRequest
		_, err := b.Validate()
		require.NoError(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		name := "   board name   "
		b := pulpeHttp.BoardUpdateRequest{
			Name: &name,
		}
		_, err := b.Validate()
		require.NoError(t, err)
		require.Equal(t, "board name", *b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := pulpeHttp.BoardUpdateRequest{
			Name: &spaces,
		}
		_, err := b.Validate()
		require.Error(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		b := pulpeHttp.BoardUpdateRequest{
			Name: &emptyName,
		}
		_, err := b.Validate()
		require.Error(t, err)
	})
}

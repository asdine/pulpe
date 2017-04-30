package api_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/http/api"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestBoardHandler_Boards(t *testing.T) {
	t.Run("OK", testBoardHandler_Boards_OK)
	t.Run("Internal error", testBoardHandler_Boards_InternalError)
	t.Run("Auth failed", testBoardHandler_Boards_AuthenticationFailed)
}

func testBoardHandler_Boards_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardsFn = func() ([]*pulpe.Board, error) {
		return []*pulpe.Board{
			&pulpe.Board{ID: "id", Name: "name", Slug: "slug", CreatedAt: mock.Now, UpdatedAt: &mock.Now, Owner: &pulpe.User{ID: "123"}},
		}, nil
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/user/boards", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `[
    {
			"id": "id",
			"name": "name",
			"slug": "slug",
			"owner": {"fullName":"", "login":"", "email":"", "id":"123", "createdAt":"0001-01-01T00:00:00Z"},
      "createdAt": `+string(date)+`,
      "updatedAt": `+string(date)+`
	  }
  ]`, w.Body.String())
}

func testBoardHandler_Boards_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardsFn = func() ([]*pulpe.Board, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/user/boards", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardsInvoked)
}

func testBoardHandler_Boards_AuthenticationFailed(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.BoardService.BoardsFn = func() ([]*pulpe.Board, error) {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/user/boards", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.True(t, c.BoardService.BoardsInvoked)
}

func TestBoardHandler_CreateBoard(t *testing.T) {
	t.Run("OK", testBoardHandler_CreateBoard_OK)
	t.Run("ErrInvalidJSON", testBoardHandler_CreateBoard_ErrInvalidJSON)
	t.Run("ValidationError", testBoardHandler_CreateBoard_ValidationError)
	t.Run("ErrInternal", testBoardHandler_CreateBoard_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
	t.Run("Authfailed", testBoardHandler_CreateBoard_WithResponse(t, http.StatusUnauthorized, pulpe.ErrUserAuthenticationFailed))
}

func testBoardHandler_CreateBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.CreateBoardFn = func(c *pulpe.BoardCreation) (*pulpe.Board, error) {
		require.Equal(t, "name", c.Name)

		return &pulpe.Board{
			ID:        "123",
			CreatedAt: mock.Now,
			Name:      c.Name,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/user/boards", bytes.NewReader([]byte(`{
    "name": "name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.CreateBoardInvoked)
}

func testBoardHandler_CreateBoard_ErrInvalidJSON(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/user/boards", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testBoardHandler_CreateBoard_ValidationError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/user/boards", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "validation error", "fields": {"name": ["non zero value required"]}}`, w.Body.String())
}

func testBoardHandler_CreateBoard_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := newHandler(c)

		// Mock service.
		c.BoardService.CreateBoardFn = func(Board *pulpe.BoardCreation) (*pulpe.Board, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/user/boards", bytes.NewReader([]byte(`{"name": "name"}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
		require.True(t, c.BoardService.CreateBoardInvoked)
	}
}

func TestBoardHandler_Board(t *testing.T) {
	t.Run("OK", testBoardHandler_Board_OK)
	t.Run("Not found", testBoardHandler_Board_NotFound)
	t.Run("Internal error", testBoardHandler_Board_InternalError)
	t.Run("Auth failed", testBoardHandler_Board_AuthenticationFailed)
}

func testBoardHandler_Board_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardByOwnerAndSlugFn = func(owner, slug string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
		require.Equal(t, "user", owner)
		require.Equal(t, "XXX", slug)
		require.Len(t, options, 2)
		return &pulpe.Board{
			ID:    "XXX",
			Owner: &pulpe.User{ID: "123"},
		}, nil
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/boards/user/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.BoardByOwnerAndSlugInvoked)
}

func testBoardHandler_Board_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardByOwnerAndSlugFn = func(owner, slug string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/boards/user/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.BoardService.BoardByOwnerAndSlugInvoked)
}

func testBoardHandler_Board_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardByOwnerAndSlugFn = func(owner, slug string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/boards/user/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.BoardByOwnerAndSlugInvoked)
}

func testBoardHandler_Board_AuthenticationFailed(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.BoardByOwnerAndSlugFn = func(owner, slug string, options ...pulpe.BoardGetOption) (*pulpe.Board, error) {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/boards/user/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.True(t, c.BoardService.BoardByOwnerAndSlugInvoked)
}

func TestBoardHandler_DeleteBoard(t *testing.T) {
	t.Run("OK", testBoardHandler_DeleteBoard_OK)
	t.Run("Not found", testBoardHandler_DeleteBoard_NotFound)
	t.Run("Internal error on delete board", testBoardHandler_DeleteBoard_InternalErrorOnDeleteBoard)
}

func testBoardHandler_DeleteBoard_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	byBoardID := func(id string) error {
		require.Equal(t, "XXX", id)
		return nil
	}

	c.BoardService.DeleteBoardFn = byBoardID

	// Retrieve Board.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
}

func testBoardHandler_DeleteBoard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.BoardService.DeleteBoardFn = func(id string) error {
		return pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
}

func testBoardHandler_DeleteBoard_InternalErrorOnDeleteBoard(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.BoardService.DeleteBoardFn = func(id string) error {
		return errors.New("unexpected error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/boards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.DeleteBoardInvoked)
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
	h := newHandler(c)

	// Mock service.
	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)

		return new(pulpe.Board), nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.NotZero(t, w.Body.Len())
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_ErrInvalidJSON(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/boards/XXX", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testBoardHandler_UpdateBoard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_ValidationError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/boards/XXX", bytes.NewReader([]byte(`{
    "name": "       "
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, c.BoardService.UpdateBoardInvoked)
}

func testBoardHandler_UpdateBoard_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.BoardService.UpdateBoardFn = func(id string, u *pulpe.BoardUpdate) (*pulpe.Board, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/boards/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.BoardService.UpdateBoardInvoked)
}

func TestBoardCreateRequest_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var b api.BoardCreateRequest
		_, err := b.Validate()
		require.Error(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		b := api.BoardCreateRequest{
			Name: "    board name   ",
		}
		_, err := b.Validate()
		require.NoError(t, err)
		require.Equal(t, "board name", b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := api.BoardCreateRequest{
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
		var b api.BoardUpdateRequest
		_, err := b.Validate()
		require.NoError(t, err)
	})

	t.Run("NameOnly", func(t *testing.T) {
		name := "   board name   "
		b := api.BoardUpdateRequest{
			Name: &name,
		}
		_, err := b.Validate()
		require.NoError(t, err)
		require.Equal(t, "board name", *b.Name)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		b := api.BoardUpdateRequest{
			Name: &spaces,
		}
		_, err := b.Validate()
		require.Error(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		b := api.BoardUpdateRequest{
			Name: &emptyName,
		}
		_, err := b.Validate()
		require.Error(t, err)
	})
}

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

func TestListHandler_CreateList(t *testing.T) {
	t.Run("OK", testListHandler_CreateList_OK)
	t.Run("ErrInvalidJSON", testListHandler_CreateList_ErrInvalidJSON)
	t.Run("ErrValidation", testListHandler_CreateList_ErrValidation)
	t.Run("NotFound", testListHandler_CreateList_BoardNotFound)
	t.Run("ErrInternal", testListHandler_CreateList_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testListHandler_CreateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.BoardFn = func(selector string) (*pulpe.Board, error) {
		require.Equal(t, "the-board", selector)
		return &pulpe.Board{
			ID: "456",
		}, nil
	}

	c.ListService.CreateListFn = func(list *pulpe.ListCreate) (*pulpe.List, error) {
		require.Equal(t, &pulpe.ListCreate{
			BoardID: "456",
			Name:    "Name",
		}, list)

		return &pulpe.List{
			ID:        "123",
			BoardID:   "456",
			Name:      "Name",
			Slug:      "slug",
			CreatedAt: mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards/the-board/lists", bytes.NewReader([]byte(`{
		"name": "Name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
    "boardID": "456",
		"name": "Name",
		"slug": "slug",
		"createdAt": `+string(date)+`
  }`, w.Body.String())
	require.True(t, c.ListService.CreateListInvoked)
	require.True(t, c.BoardService.BoardInvoked)
}

func testListHandler_CreateList_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards/the-board/lists", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testListHandler_CreateList_ErrValidation(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards/the-board/lists", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "validation error", "fields": {"name": ["non zero value required"]}}`, w.Body.String())
}

func testListHandler_CreateList_BoardNotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.BoardService.BoardFn = func(selector string) (*pulpe.Board, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/boards/the-board/lists", bytes.NewReader([]byte(`{
		"name": "Name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.False(t, c.ListService.CreateListInvoked)
}

func testListHandler_CreateList_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		// Mock service.
		c.ListService.CreateListFn = func(list *pulpe.ListCreate) (*pulpe.List, error) {
			return nil, err
		}

		c.BoardService.BoardFn = func(selector string) (*pulpe.Board, error) {
			require.Equal(t, "the-board", selector)
			return new(pulpe.Board), nil
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/boards/the-board/lists", bytes.NewReader([]byte(`{"name": "name"}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func TestListHandler_DeleteList(t *testing.T) {
	t.Run("OK", testListHandler_DeleteList_OK)
	t.Run("Not found", testListHandler_DeleteList_NotFound)
	t.Run("Internal error on delete list", testListHandler_DeleteList_InternalErrorOnDeleteList)
	t.Run("Internal error on delete cards by list id", testListHandler_DeleteList_InternalErrorOnDeleteCardsByListID)
}

func testListHandler_DeleteList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	c.CardService.DeleteCardsByListIDFn = func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
	require.True(t, c.CardService.DeleteCardsByListIDInvoked)
}

func testListHandler_DeleteList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return pulpe.ErrListNotFound
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
	require.False(t, c.CardService.DeleteCardsByListIDInvoked)
}

func testListHandler_DeleteList_InternalErrorOnDeleteList(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return errors.New("unexpected error")
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
	require.False(t, c.CardService.DeleteCardsByListIDInvoked)
}

func testListHandler_DeleteList_InternalErrorOnDeleteCardsByListID(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return nil
	}

	c.CardService.DeleteCardsByListIDFn = func(id string) error {
		return errors.New("unexpected error")
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
	require.True(t, c.CardService.DeleteCardsByListIDInvoked)
}

func TestListHandler_UpdateList(t *testing.T) {
	t.Run("OK", testListHandler_UpdateList_OK)
	t.Run("ErrInvalidJSON", testListHandler_UpdateList_ErrInvalidJSON)
	t.Run("Not found", testListHandler_UpdateList_NotFound)
	t.Run("Internal error", testListHandler_UpdateList_InternalError)
	t.Run("Validation error", testListHandler_UpdateList_ValidationError)
}

func testListHandler_UpdateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)
		return &pulpe.List{
			ID:        "XXX",
			Name:      *u.Name,
			Slug:      "new-name",
			CreatedAt: mock.Now,
			UpdatedAt: &mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "XXX",
    "name": "new name",
		"slug": "new-name",
    "boardID": "",
		"createdAt": `+string(date)+`,
		"updatedAt": `+string(date)+`
  }`, w.Body.String())
}

func testListHandler_UpdateList_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testListHandler_UpdateList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, pulpe.ErrListNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func testListHandler_UpdateList_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func testListHandler_UpdateList_ValidationError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": ""
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, c.ListService.UpdateListInvoked)
}

func TestListCreateRequest_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var l pulpeHttp.ListCreateRequest
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := pulpeHttp.ListCreateRequest{
			Name: "      ",
		}
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		l := pulpeHttp.ListCreateRequest{
			Name: "list name",
		}
		_, err := l.Validate()
		require.NoError(t, err)
	})
}

func TestListUpdate_Validate(t *testing.T) {
	name := "list name"
	emptyName := ""
	spaces := "    "

	t.Run("Empty", func(t *testing.T) {
		var l pulpeHttp.ListUpdateRequest
		_, err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("ValidName", func(t *testing.T) {
		l := pulpeHttp.ListUpdateRequest{
			Name: &name,
		}
		_, err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := pulpeHttp.ListUpdateRequest{
			Name: &spaces,
		}
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		l := pulpeHttp.ListUpdateRequest{
			Name: &emptyName,
		}
		_, err := l.Validate()
		require.Error(t, err)
	})
}

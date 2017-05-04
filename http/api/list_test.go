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

func TestListHandler_CreateList(t *testing.T) {
	t.Run("OK", testListHandler_CreateList_OK)
	t.Run("ErrInvalidJSON", testListHandler_CreateList_ErrInvalidJSON)
	t.Run("ErrValidation", testListHandler_CreateList_ErrValidation)
	t.Run("NotFound", testListHandler_CreateList_BoardNotFound)
	t.Run("ErrInternal", testListHandler_CreateList_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
	t.Run("AuthFailed", testListHandler_CreateList_WithResponse(t, http.StatusUnauthorized, pulpe.ErrUserAuthenticationFailed))
}

func testListHandler_CreateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.CreateListFn = func(boardID string, list *pulpe.ListCreation) (*pulpe.List, error) {
		require.Equal(t, &pulpe.ListCreation{
			Name:     "Name",
			Position: 123.45,
		}, list)

		return &pulpe.List{
			ID:        "123",
			BoardID:   boardID,
			Name:      "Name",
			Position:  123.45,
			Slug:      "slug",
			OwnerID:   "456",
			CreatedAt: mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/boards/XXX/lists", bytes.NewReader([]byte(`{
		"name": "Name",
		"position": 123.45
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
    "boardID": "XXX",
		"position": 123.45,
		"name": "Name",
		"slug": "slug",
		"ownerID": "456",
		"createdAt": `+string(date)+`
  }`, w.Body.String())
	require.True(t, c.ListService.CreateListInvoked)
}

func testListHandler_CreateList_ErrInvalidJSON(t *testing.T) {
	h := newHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/boards/the-board/lists", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testListHandler_CreateList_ErrValidation(t *testing.T) {
	h := newHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/boards/the-board/lists", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "validation error", "fields": {"name": ["non zero value required"]}}`, w.Body.String())
}

func testListHandler_CreateList_BoardNotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.CreateListFn = func(boardID string, list *pulpe.ListCreation) (*pulpe.List, error) {
		return nil, pulpe.ErrBoardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/boards/the-board/lists", bytes.NewReader([]byte(`{
		"name": "Name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.ListService.CreateListInvoked)
}

func testListHandler_CreateList_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := newHandler(c)

		// Mock service.
		c.ListService.CreateListFn = func(boardID string, list *pulpe.ListCreation) (*pulpe.List, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/boards/the-board/lists", bytes.NewReader([]byte(`{"name": "name"}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func TestListHandler_DeleteList(t *testing.T) {
	t.Run("OK", testListHandler_DeleteList_OK)
	t.Run("Not found", testListHandler_DeleteList_NotFound)
	t.Run("Internal error on delete list", testListHandler_DeleteList_InternalErrorOnDeleteList)
	t.Run("Auth failed on delete list", testListHandler_DeleteList_AuthErrorOnDeleteList)
}

func testListHandler_DeleteList_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.DeleteListFn = func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
}

func testListHandler_DeleteList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return pulpe.ErrListNotFound
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
}

func testListHandler_DeleteList_InternalErrorOnDeleteList(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return errors.New("unexpected error")
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
}

func testListHandler_DeleteList_AuthErrorOnDeleteList(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id string) error {
		return pulpe.ErrUserAuthenticationFailed
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/api/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.True(t, c.ListService.DeleteListInvoked)
}

func TestListHandler_UpdateList(t *testing.T) {
	t.Run("OK", testListHandler_UpdateList_OK)
	t.Run("ErrInvalidJSON", testListHandler_UpdateList_ErrInvalidJSON)
	t.Run("Not found", testListHandler_UpdateList_NotFound)
	t.Run("Internal error", testListHandler_UpdateList_InternalError)
	t.Run("Validation error", testListHandler_UpdateList_ValidationError)
	t.Run("Auth failed", testListHandler_UpdateList_AuthenticationFailed)
}

func testListHandler_UpdateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	// Mock service.
	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		require.Equal(t, "new name", *u.Name)
		require.Equal(t, 123.45, *u.Position)
		require.Equal(t, "new name", *u.Name)
		return &pulpe.List{
			ID:        "XXX",
			Name:      *u.Name,
			Position:  *u.Position,
			Slug:      "new-name",
			BoardID:   "XXX",
			OwnerID:   "456",
			CreatedAt: mock.Now,
			UpdatedAt: &mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name",
		"position": 123.45
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "XXX",
    "name": "new name",
		"position": 123.45,
		"slug": "new-name",
    "boardID": "XXX",
		"ownerID": "456",
		"createdAt": `+string(date)+`,
		"updatedAt": `+string(date)+`
  }`, w.Body.String())
}

func testListHandler_UpdateList_ErrInvalidJSON(t *testing.T) {
	h := newHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testListHandler_UpdateList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, pulpe.ErrListNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func testListHandler_UpdateList_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func testListHandler_UpdateList_ValidationError(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "name": ""
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, c.ListService.UpdateListInvoked)
}

func testListHandler_UpdateList_AuthenticationFailed(t *testing.T) {
	c := mock.NewClient()
	h := newHandler(c)

	c.ListService.UpdateListFn = func(id string, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/api/lists/XXX", bytes.NewReader([]byte(`{
    "name": ""
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.False(t, c.ListService.UpdateListInvoked)
}

func TestListCreateRequest_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var l api.ListCreateRequest
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := api.ListCreateRequest{
			Name: "      ",
		}
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		l := api.ListCreateRequest{
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
		var l api.ListUpdateRequest
		_, err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("ValidName", func(t *testing.T) {
		l := api.ListUpdateRequest{
			Name: &name,
		}
		_, err := l.Validate()
		require.NoError(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		l := api.ListUpdateRequest{
			Name: &spaces,
		}
		_, err := l.Validate()
		require.Error(t, err)
	})

	t.Run("EmptyName", func(t *testing.T) {
		l := api.ListUpdateRequest{
			Name: &emptyName,
		}
		_, err := l.Validate()
		require.Error(t, err)
	})
}

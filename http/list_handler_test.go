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
	t.Run("ErrListIDRequired", testListHandler_CreateList_WithResponse(t, http.StatusBadRequest, pulpe.ErrListIDRequired))
	t.Run("ErrListBoardIDRequired", testListHandler_CreateList_WithResponse(t, http.StatusBadRequest, pulpe.ErrListBoardIDRequired))
	t.Run("ErrListExists", testListHandler_CreateList_WithResponse(t, http.StatusConflict, pulpe.ErrListExists))
	t.Run("ErrInternal", testListHandler_CreateList_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testListHandler_CreateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.CreateListFn = func(list *pulpe.ListCreate) (*pulpe.List, error) {
		require.Equal(t, &pulpe.ListCreate{
			BoardID: "456",
			Name:    "Name",
		}, list)

		return &pulpe.List{
			ID:        "123",
			BoardID:   "456",
			Name:      "Name",
			CreatedAt: mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists", bytes.NewReader([]byte(`{
    "boardID": "456",
		"name": "Name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)

	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
    "boardID": "456",
		"name": "Name",
		"createdAt": `+string(date)+`
  }`, w.Body.String())
}

func testListHandler_CreateList_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid json"}`, w.Body.String())
}

func testListHandler_CreateList_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		// Mock service.
		c.ListService.CreateListFn = func(list *pulpe.ListCreate) (*pulpe.List, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/lists", bytes.NewReader([]byte(`{}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func TestListHandler_DeleteList(t *testing.T) {
	t.Run("OK", testListHandler_DeleteList_OK)
	t.Run("Not found", testListHandler_DeleteList_NotFound)
	t.Run("Internal error", testListHandler_DeleteList_InternalError)
}

func testListHandler_DeleteList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id pulpe.ListID) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func testListHandler_DeleteList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id pulpe.ListID) error {
		return pulpe.ErrListNotFound
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, `{}`, w.Body.String())
}

func testListHandler_DeleteList_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.DeleteListFn = func(id pulpe.ListID) error {
		return errors.New("unexpected error")
	}

	// Retrieve List.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/lists/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListHandler_UpdateList(t *testing.T) {
	t.Run("OK", testListHandler_UpdateList_OK)
	t.Run("ErrInvalidJSON", testListHandler_UpdateList_ErrInvalidJSON)
	t.Run("Not found", testListHandler_UpdateList_NotFound)
	t.Run("Internal error", testListHandler_UpdateList_InternalError)
}

func testListHandler_UpdateList_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.ListService.UpdateListFn = func(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)
		return &pulpe.List{
			ID:        "XXX",
			Name:      *u.Name,
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
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "XXX",
    "name": "new name",
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
	require.JSONEq(t, `{"err": "invalid json"}`, w.Body.String())
}

func testListHandler_UpdateList_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.UpdateListFn = func(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, pulpe.ErrListNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, `{}`, w.Body.String())
}

func testListHandler_UpdateList_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.UpdateListFn = func(id pulpe.ListID, u *pulpe.ListUpdate) (*pulpe.List, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/lists/XXX", bytes.NewReader([]byte(`{
    "name": "new name"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

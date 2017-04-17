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

func TestCardHandler_CreateCard(t *testing.T) {
	t.Run("OK", testCardHandler_CreateCard_OK)
	t.Run("ErrInvalidJSON", testCardHandler_CreateCard_ErrInvalidJSON)
	t.Run("ErrValidation", testCardHandler_CreateCard_ErrValidation)
	t.Run("NotFound", testCardHandler_CreateCard_ListNotFound)
	t.Run("ErrInternal", testCardHandler_CreateCard_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testCardHandler_CreateCard_OK(t *testing.T) {
	c := mock.NewClient()

	c.ListService.ListFn = func(id string) (*pulpe.List, error) {
		require.Equal(t, "456", string(id))
		return &pulpe.List{
			ID:      "456",
			BoardID: "789",
		}, nil
	}

	c.CardService.CreateCardFn = func(c *pulpe.CardCreation) (*pulpe.Card, error) {
		require.Equal(t, &pulpe.CardCreation{
			ListID:      "456",
			BoardID:     "789",
			Name:        "name",
			Description: "description",
			Position:    1,
		}, c)

		return &pulpe.Card{
			ID:          "123",
			CreatedAt:   mock.Now,
			Slug:        "slug",
			ListID:      c.ListID,
			BoardID:     c.BoardID,
			Name:        c.Name,
			Description: c.Description,
			Position:    c.Position,
		}, nil
	}

	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists/456/cards", bytes.NewReader([]byte(`{
    "name": "name",
    "description": "description",
		"position": 1
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
    "listID": "456",
    "boardID": "789",
    "name": "name",
		"slug": "slug",
    "description": "description",
		"position": 1,
		"createdAt": `+string(date)+`
  }`, w.Body.String())
}

func testCardHandler_CreateCard_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists/456/cards", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testCardHandler_CreateCard_ListNotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.ListService.ListFn = func(id string) (*pulpe.List, error) {
		return nil, pulpe.ErrListNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists/abcd/cards", bytes.NewReader([]byte(`{
		"name": "Name",
		"position": 10
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
	require.False(t, c.CardService.CreateCardInvoked)
}

func testCardHandler_CreateCard_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		c.CardService.CreateCardFn = func(card *pulpe.CardCreation) (*pulpe.Card, error) {
			return nil, err
		}

		c.ListService.ListFn = func(id string) (*pulpe.List, error) {
			require.Equal(t, "456", string(id))
			return new(pulpe.List), nil
		}

		c.BoardService.BoardFn = func(id string) (*pulpe.Board, error) {
			require.Equal(t, "789", string(id))
			return new(pulpe.Board), nil
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/lists/456/cards", bytes.NewReader([]byte(`{
			"listID": "456",
			"name": "name",
			"description": "description",
			"position": 1
		}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func testCardHandler_CreateCard_ErrValidation(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/lists/456/cards", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCardHandler_Card(t *testing.T) {
	t.Run("OK", testCardHandler_Card_OK)
	t.Run("Not found", testCardHandler_Card_NotFound)
	t.Run("Internal error", testCardHandler_Card_InternalError)
}

func testCardHandler_Card_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.CardFn = func(id string) (*pulpe.Card, error) {
		require.Equal(t, "XXX", string(id))
		return &pulpe.Card{ID: id, Name: "name", Description: "description", Position: 2 << 3, ListID: "YYY", BoardID: "ZZZ", CreatedAt: mock.Now, UpdatedAt: &mock.Now, Slug: "slug"}, nil
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "XXX",
		"listID": "YYY",
		"boardID": "ZZZ",
    "name": "name",
		"slug": "slug",
    "description": "description",
		"position": 16,
    "createdAt": `+string(date)+`,
    "updatedAt": `+string(date)+`
		}`, w.Body.String())
}

func testCardHandler_Card_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.CardFn = func(id string) (*pulpe.Card, error) {
		return nil, pulpe.ErrCardNotFound
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func testCardHandler_Card_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.CardFn = func(id string) (*pulpe.Card, error) {
		return nil, errors.New("unexpected error")
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCardHandler_DeleteCard(t *testing.T) {
	t.Run("OK", testCardHandler_DeleteCard_OK)
	t.Run("Not found", testCardHandler_DeleteCard_NotFound)
	t.Run("Internal error", testCardHandler_DeleteCard_InternalError)
}

func testCardHandler_DeleteCard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.DeleteCardFn = func(id string) error {
		require.Equal(t, "XXX", string(id))
		return nil
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNoContent, w.Code)
}

func testCardHandler_DeleteCard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.DeleteCardFn = func(id string) error {
		return pulpe.ErrCardNotFound
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func testCardHandler_DeleteCard_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.DeleteCardFn = func(id string) error {
		return errors.New("unexpected error")
	}

	// Retrieve Card.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/v1/cards/XXX", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCardHandler_UpdateCard(t *testing.T) {
	t.Run("OK", testCardHandler_UpdateCard_OK)
	t.Run("ErrInvalidJSON", testCardHandler_UpdateCard_ErrInvalidJSON)
	t.Run("Not found", testCardHandler_UpdateCard_NotFound)
	t.Run("Validation error", testCardHandler_UpdateCard_ErrValidation)
	t.Run("Internal error", testCardHandler_UpdateCard_InternalError)
}

func testCardHandler_UpdateCard_OK(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	// Mock service.
	c.CardService.UpdateCardFn = func(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
		require.Equal(t, "XXX", string(id))
		require.NotNil(t, u.Name)
		require.Equal(t, "new name", *u.Name)
		require.NotNil(t, u.Description)
		require.Zero(t, *u.Description)
		require.NotNil(t, u.Position)
		require.Zero(t, *u.Position)
		return &pulpe.Card{
			ID:          "XXX",
			Name:        *u.Name,
			Description: *u.Description,
			Position:    *u.Position,
			CreatedAt:   mock.Now,
			UpdatedAt:   &mock.Now,
		}, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/cards/XXX", bytes.NewReader([]byte(`{
    "name": "new name",
    "description": "",
    "position": 0
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "XXX",
    "name": "new name",
		"slug": "",
		"description": "",
		"position": 0,
    "listID": "",
    "boardID": "",
		"createdAt": `+string(date)+`,
		"updatedAt": `+string(date)+`
  }`, w.Body.String())
}

func testCardHandler_UpdateCard_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/cards/XXX", bytes.NewReader([]byte(`{
    "id": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testCardHandler_UpdateCard_ErrValidation(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/cards/XXX", bytes.NewReader([]byte(`{
		"name": ""
	}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func testCardHandler_UpdateCard_NotFound(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.CardService.UpdateCardFn = func(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
		return nil, pulpe.ErrCardNotFound
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/cards/XXX", bytes.NewReader([]byte(`{
    "name": "new name",
    "description": ""
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func testCardHandler_UpdateCard_InternalError(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.CardService.UpdateCardFn = func(id string, u *pulpe.CardUpdate) (*pulpe.Card, error) {
		return nil, errors.New("internal error")
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PATCH", "/v1/cards/XXX", bytes.NewReader([]byte(`{
    "name": "new name",
    "description": ""
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCardCreateRequest_Validate(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		var cc pulpeHttp.CardCreateRequest
		_, err := cc.Validate()
		require.Error(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		cc := pulpeHttp.CardCreateRequest{
			Name: "      ",
		}
		_, err := cc.Validate()
		require.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		cc := pulpeHttp.CardCreateRequest{
			Name: "Card name",
		}
		_, err := cc.Validate()
		require.NoError(t, err)
	})

	t.Run("NegativePosition", func(t *testing.T) {
		cc := pulpeHttp.CardCreateRequest{
			Name:     "Card name",
			Position: -10.0,
		}
		_, err := cc.Validate()
		require.Error(t, err)
	})
}

func TestCardUpdate_Validate(t *testing.T) {
	name := "Card name"
	emptyName := ""
	spaces := "    "
	zeroFloat := 0.0
	negativeFloat := -10.0
	positiveFloat := 10.0

	t.Run("Empty", func(t *testing.T) {
		var cc pulpeHttp.CardUpdateRequest
		_, err := cc.Validate()
		require.NoError(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		cc := pulpeHttp.CardUpdateRequest{
			Name:        &name,
			Description: &name,
			Position:    &positiveFloat,
		}
		_, err := cc.Validate()
		require.NoError(t, err)
	})

	t.Run("SpaceOnly", func(t *testing.T) {
		cc := pulpeHttp.CardUpdateRequest{
			Name:        &spaces,
			Description: &spaces,
		}
		_, err := cc.Validate()
		require.Error(t, err)
	})

	t.Run("ZeroValues", func(t *testing.T) {
		cc := pulpeHttp.CardUpdateRequest{
			Name:        &emptyName,
			Description: &emptyName,
			Position:    &zeroFloat,
		}
		_, err := cc.Validate()
		require.Error(t, err)
	})

	t.Run("NegativePosition", func(t *testing.T) {
		cc := pulpeHttp.CardUpdateRequest{
			Name:        &name,
			Description: &name,
			Position:    &negativeFloat,
		}
		_, err := cc.Validate()
		require.Error(t, err)
	})
}

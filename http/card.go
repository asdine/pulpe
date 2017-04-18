package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/validation"
	"github.com/julienschmidt/httprouter"
)

// NewCardHandler returns a new instance of CardHandler.
func NewCardHandler(router *httprouter.Router, c pulpe.Client) *CardHandler {
	h := CardHandler{
		Router: router,
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/v1/lists/:listID/cards", h.handlePostCard)
	h.GET("/v1/cards/:id", h.handleGetCard)
	h.DELETE("/v1/cards/:id", h.handleDeleteCard)
	h.PATCH("/v1/cards/:id", h.handlePatchCard)
	return &h
}

// CardHandler represents an HTTP API handler for cards.
type CardHandler struct {
	*httprouter.Router

	Client pulpe.Client
	Logger *log.Logger
}

// handlePostCard handles requests to create a new card.
func (h *CardHandler) handlePostCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req CardCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	cc, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	// fetch list
	listID := ps.ByName("listID")
	list, err := session.ListService().List(listID)
	if err != nil {
		if err == pulpe.ErrListNotFound {
			http.NotFound(w, r)
		} else {
			Error(w, err, http.StatusInternalServerError, h.Logger)
		}
		return
	}

	cc.ListID = list.ID
	cc.BoardID = list.BoardID

	card, err := session.CardService().CreateCard(cc)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusCreated, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleGetCard handles requests to fetch a single card.
func (h *CardHandler) handleGetCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()
	defer session.Close()

	card, err := session.CardService().Card(id)
	if err != nil {
		if err == pulpe.ErrCardNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, card, http.StatusOK, h.Logger)
}

// handleDeleteCard handles requests to delete a single card.
func (h *CardHandler) handleDeleteCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()
	defer session.Close()

	err := session.CardService().DeleteCard(id)
	if err != nil {
		if err == pulpe.ErrCardNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchCard handles requests to update a card.
func (h *CardHandler) handlePatchCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req CardUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	cu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	card, err := session.CardService().UpdateCard(id, cu)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.Logger)
	case pulpe.ErrCardNotFound:
		http.NotFound(w, r)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// CardCreateRequest is the payload sent to create a card.
type CardCreateRequest struct {
	Name        string  `json:"name" valid:"required,stringlength(1|64)"`
	Description string  `json:"description" valid:"stringlength(1|100000)"`
	Position    float64 `json:"position"`
}

// Validate list creation payload.
func (c *CardCreateRequest) Validate() (*pulpe.CardCreation, error) {
	c.Name = strings.TrimSpace(c.Name)
	c.Description = strings.TrimSpace(c.Description)

	verr := validation.Validate(c)

	// validate position
	if c.Position < 0 {
		verr = validation.AddError(verr, "position", errors.New("position should be greater than zero"))
	}

	if verr != nil {
		return nil, verr
	}

	return &pulpe.CardCreation{
		Name:        c.Name,
		Description: c.Description,
		Position:    c.Position,
	}, nil
}

// CardUpdateRequest is the payload sent to update a card.
type CardUpdateRequest struct {
	Name        *string  `json:"name" valid:"stringlength(1|64)"`
	Description *string  `json:"description" valid:"stringlength(1|100000)"`
	Position    *float64 `json:"position"`
}

// Validate card update payload.
func (c *CardUpdateRequest) Validate() (*pulpe.CardUpdate, error) {
	if c.Name != nil {
		*c.Name = strings.TrimSpace(*c.Name)
	}

	if c.Description != nil {
		*c.Description = strings.TrimSpace(*c.Description)
	}

	err := validation.Validate(c)
	if c.Name != nil && *c.Name == "" {
		err = validation.AddError(err, "name", errors.New("name should not be empty"))
	}

	if c.Position != nil && *c.Position < 0 {
		err = validation.AddError(err, "position", errors.New("position should be greater than zero"))
	}

	if err != nil {
		return nil, err
	}

	return &pulpe.CardUpdate{
		Name:        c.Name,
		Description: c.Description,
		Position:    c.Position,
	}, nil
}

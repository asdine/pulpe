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

// registerCardHandler register the cardHandler routes.
func registerCardHandler(router *httprouter.Router, c *client) {
	h := cardHandler{
		client: c,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	router.POST("/api/lists/:listID/cards", h.handlePostCard)
	router.GET("/api/cards/:id", h.handleGetCard)
	router.DELETE("/api/cards/:id", h.handleDeleteCard)
	router.PATCH("/api/cards/:id", h.handlePatchCard)
}

// cardHandler represents an HTTP API handler for cards.
type cardHandler struct {
	client *client
	logger *log.Logger
}

// handlePostCard handles requests to create a new card.
func (h *cardHandler) handlePostCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req CardCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	cc, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	listID := ps.ByName("listID")

	card, err := session.CardService().CreateCard(listID, cc)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusCreated, h.logger)
	case pulpe.ErrListNotFound:
		Error(w, err, http.StatusNotFound, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handleGetCard handles requests to fetch a single card.
func (h *cardHandler) handleGetCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.client.session(w, r)
	defer session.Close()

	card, err := session.CardService().Card(id)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.logger)
	case pulpe.ErrCardNotFound:
		Error(w, err, http.StatusNotFound, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handleDeleteCard handles requests to delete a single card.
func (h *cardHandler) handleDeleteCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.client.session(w, r)
	defer session.Close()

	err := session.CardService().DeleteCard(id)
	switch err {
	case nil:
		w.WriteHeader(http.StatusNoContent)
	case pulpe.ErrCardNotFound:
		Error(w, err, http.StatusNotFound, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handlePatchCard handles requests to update a card.
func (h *cardHandler) handlePatchCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req CardUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	cu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	card, err := session.CardService().UpdateCard(id, cu)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.logger)
	case pulpe.ErrCardNotFound:
		Error(w, err, http.StatusNotFound, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
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

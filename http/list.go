package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/validation"
	"github.com/julienschmidt/httprouter"
)

// registerListHandler register the listHandler routes.
func registerListHandler(router *httprouter.Router, c *client) {
	h := listHandler{
		client: c,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	router.POST("/boards/:boardID/lists", h.handlePostList)
	router.DELETE("/lists/:id", h.handleDeleteList)
	router.PATCH("/lists/:id", h.handlePatchList)
}

// listHandler represents an HTTP API handler for lists.
type listHandler struct {
	*httprouter.Router

	client *client

	logger *log.Logger
}

// handlePostList handles requests to create a new list.
func (h *listHandler) handlePostList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	boardID := ps.ByName("boardID")

	// decode request
	var req ListCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	// validate payload
	lc, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	// create the list
	list, err := session.ListService().CreateList(boardID, lc)
	switch err {
	case nil:
		encodeJSON(w, list, http.StatusCreated, h.logger)
	case pulpe.ErrBoardNotFound:
		Error(w, err, http.StatusNotFound, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handleDeleteList handles requests to delete a single list and all of its cards.
func (h *listHandler) handleDeleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.client.session(w, r)
	defer session.Close()

	err := session.ListService().DeleteList(id)
	if err != nil {
		if err == pulpe.ErrListNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchList handles requests to update a list.
func (h *listHandler) handlePatchList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req ListUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	lu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	card, err := session.ListService().UpdateList(id, lu)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.logger)
	case pulpe.ErrListNotFound:
		http.NotFound(w, r)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// ListCreateRequest is used to create a List.
type ListCreateRequest struct {
	Name string `json:"name" valid:"required,stringlength(1|64)"`
}

// Validate list creation payload.
func (l *ListCreateRequest) Validate() (*pulpe.ListCreation, error) {
	l.Name = strings.TrimSpace(l.Name)

	err := validation.Validate(l)
	if err != nil {
		return nil, err
	}

	return &pulpe.ListCreation{
		Name: l.Name,
	}, nil
}

// ListUpdateRequest is used to update a List.
type ListUpdateRequest struct {
	Name *string `json:"name" valid:"required,stringlength(1|64)"`
}

// Validate list update payload.
func (l *ListUpdateRequest) Validate() (*pulpe.ListUpdate, error) {
	if l.Name == nil {
		return nil, nil
	}

	*l.Name = strings.TrimSpace(*l.Name)

	err := validation.Validate(l)
	if err != nil {
		return nil, err
	}

	return &pulpe.ListUpdate{
		Name: l.Name,
	}, nil
}

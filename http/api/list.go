package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/blankrobot/pulpe"
	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/validation"
	"github.com/julienschmidt/httprouter"
)

// registerListHandler register the listHandler routes.
func registerListHandler(router *httprouter.Router, c pulpeHttp.Connector) {
	h := listHandler{
		connect: c,
		logger:  log.New(os.Stderr, "", log.LstdFlags),
	}

	router.POST("/api/boards/:boardID/lists", h.handlePostList)
	router.DELETE("/api/lists/:id", h.handleDeleteList)
	router.PATCH("/api/lists/:id", h.handlePatchList)
}

// listHandler represents an HTTP API handler for lists.
type listHandler struct {
	connect pulpeHttp.Connector
	logger  *log.Logger
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

	session := h.connect(r)
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

	session := h.connect(r)
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

	session := h.connect(r)
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
	Name     string  `json:"name" valid:"required,stringlength(1|64)"`
	Position float64 `json:"position"`
}

// Validate list creation payload.
func (l *ListCreateRequest) Validate() (*pulpe.ListCreation, error) {
	l.Name = strings.TrimSpace(l.Name)

	err := validation.Validate(l)
	if err != nil {
		return nil, err
	}

	return &pulpe.ListCreation{
		Name:     l.Name,
		Position: l.Position,
	}, nil
}

// ListUpdateRequest is used to update a List.
type ListUpdateRequest struct {
	Name     *string  `json:"name" valid:"stringlength(1|64)"`
	Position *float64 `json:"position"`
}

// Validate list update payload.
func (l *ListUpdateRequest) Validate() (*pulpe.ListUpdate, error) {
	if l.Name != nil {
		*l.Name = strings.TrimSpace(*l.Name)
	}

	err := validation.Validate(l)

	if l.Name != nil && len(*l.Name) == 0 {
		err = validation.AddError(err, "name", errors.New("should not be empty"))
	}

	if err != nil {
		return nil, err
	}

	return &pulpe.ListUpdate{
		Name:     l.Name,
		Position: l.Position,
	}, nil
}

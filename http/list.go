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

// NewListHandler returns a new instance of ListHandler.
func NewListHandler(router *httprouter.Router, c pulpe.Client) *ListHandler {
	h := ListHandler{
		Router: router,
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/v1/boards/:board/lists", h.handlePostList)
	h.DELETE("/v1/lists/:id", h.handleDeleteList)
	h.PATCH("/v1/lists/:id", h.handlePatchList)
	return &h
}

// ListHandler represents an HTTP API handler for lists.
type ListHandler struct {
	*httprouter.Router

	Client pulpe.Client

	Logger *log.Logger
}

// handlePostList handles requests to create a new list.
func (h *ListHandler) handlePostList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode request
	var req ListCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	// validate payload
	lc, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	boardSelector := ps.ByName("board")

	// fetch board
	board, err := session.BoardService().Board(boardSelector)
	if err != nil {
		if err == pulpe.ErrBoardNotFound {
			http.NotFound(w, r)
		} else {
			Error(w, err, http.StatusInternalServerError, h.Logger)
		}
		return
	}

	// set the boardID to the ListCreate
	lc.BoardID = board.ID

	// create the list
	list, err := session.ListService().CreateList(lc)
	switch err {
	case nil:
		encodeJSON(w, list, http.StatusCreated, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleDeleteList handles requests to delete a single list and all of its cards.
func (h *ListHandler) handleDeleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()
	defer session.Close()

	err := session.ListService().DeleteList(id)
	if err != nil {
		if err == pulpe.ErrListNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	err = session.CardService().DeleteCardsByListID(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchList handles requests to update a list.
func (h *ListHandler) handlePatchList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req ListUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	lu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	card, err := session.ListService().UpdateList(id, lu)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.Logger)
	case pulpe.ErrListNotFound:
		http.NotFound(w, r)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
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

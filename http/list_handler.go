package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

// NewListHandler returns a new instance of ListHandler.
func NewListHandler(c pulpe.Client) *ListHandler {
	h := ListHandler{
		Router: httprouter.New(),
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/v1/lists", h.handlePostList)
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
func (h *ListHandler) handlePostList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request.
	var req pulpe.ListCreate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	err = req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	list, err := session.ListService().CreateList(&req)
	switch err {
	case nil:
		encodeJSON(w, list, http.StatusCreated, h.Logger)
	case pulpe.ErrListIDRequired, pulpe.ErrListBoardIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case pulpe.ErrListExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleDeleteList handles requests to delete a single list and all of its cards.
func (h *ListHandler) handleDeleteList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()
	defer session.Close()

	err := session.ListService().DeleteList(pulpe.ListID(id))
	if err != nil {
		if err == pulpe.ErrListNotFound {
			NotFound(w)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	err = session.CardService().DeleteCardsByListID(pulpe.ListID(id))
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchList handles requests to update a list.
func (h *ListHandler) handlePatchList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req pulpe.ListUpdate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	err = req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	card, err := session.ListService().UpdateList(pulpe.ListID(id), &req)
	switch err {
	case nil:
		encodeJSON(w, card, http.StatusOK, h.Logger)
	case pulpe.ErrListNotFound:
		NotFound(w)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

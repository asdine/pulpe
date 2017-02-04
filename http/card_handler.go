package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

// NewCardHandler returns a new instance of CardHandler.
func NewCardHandler(c pulpe.Client) *CardHandler {
	h := CardHandler{
		Router: httprouter.New(),
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/v1/cards", h.handlePostCard)
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
func (h *CardHandler) handlePostCard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req pulpe.CardCreate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	card, err := h.Client.Connect().CardService().CreateCard(&req)
	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		encodeJSON(w, card, h.Logger)
	case pulpe.ErrCardIDRequired, pulpe.ErrCardListIDRequired, pulpe.ErrCardBoardIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case pulpe.ErrCardExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleGetCard handles requests to fetch a single card.
func (h *CardHandler) handleGetCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	card, err := h.Client.Connect().CardService().Card(pulpe.CardID(id))
	if err != nil {
		if err == pulpe.ErrCardNotFound {
			NotFound(w)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, card, h.Logger)
}

// handleDeleteCard handles requests to delete a single card.
func (h *CardHandler) handleDeleteCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	err := h.Client.Connect().CardService().DeleteCard(pulpe.CardID(id))
	if err != nil {
		if err == pulpe.ErrCardNotFound {
			NotFound(w)
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

	var req pulpe.CardUpdate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	card, err := h.Client.Connect().CardService().UpdateCard(pulpe.CardID(id), &req)
	switch err {
	case nil:
		encodeJSON(w, card, h.Logger)
	case pulpe.ErrCardNotFound:
		NotFound(w)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

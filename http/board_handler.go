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

// NewBoardHandler returns a new instance of BoardHandler.
func NewBoardHandler(router *httprouter.Router, c pulpe.Client) *BoardHandler {
	h := BoardHandler{
		Router: router,
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.GET("/v1/boards", h.handleGetBoards)
	h.POST("/v1/boards", h.handlePostBoard)
	h.GET("/v1/boards/:board", h.handleGetBoard)
	h.DELETE("/v1/boards/:id", h.handleDeleteBoard)
	h.PATCH("/v1/boards/:id", h.handlePatchBoard)
	return &h
}

// BoardHandler represents an HTTP API handler for boards.
type BoardHandler struct {
	*httprouter.Router

	Client pulpe.Client

	Logger *log.Logger
}

// handlePostBoard handles requests to create a new board.
func (h *BoardHandler) handleGetBoards(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := h.Client.Connect()
	defer session.Close()

	boards, err := session.BoardService().Boards()
	switch err {
	case nil:
		encodeJSON(w, boards, http.StatusOK, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handlePostBoard handles requests to create a new board.
func (h *BoardHandler) handlePostBoard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req BoardCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	cr, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	board, err := session.BoardService().CreateBoard(cr)
	switch err {
	case nil:
		encodeJSON(w, board, http.StatusCreated, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// handleGetBoard handles requests to fetch a single board.
func (h *BoardHandler) handleGetBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	selector := ps.ByName("board")

	session := h.Client.Connect()
	defer session.Close()

	// Get the board
	board, err := session.BoardService().Board(selector)
	if err != nil {
		if err == pulpe.ErrBoardNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	// Get the board's lists
	board.Lists, err = session.ListService().ListsByBoard(board.ID)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	// Get the board's cards
	board.Cards, err = session.CardService().CardsByBoard(board.ID)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	encodeJSON(w, board, http.StatusOK, h.Logger)
}

// handleDeleteBoard handles requests to delete a single board and all of its content.
func (h *BoardHandler) handleDeleteBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.Client.Connect()
	defer session.Close()

	err := session.BoardService().DeleteBoard(id)
	if err != nil {
		if err == pulpe.ErrBoardNotFound {
			http.NotFound(w, r)
			return
		}

		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	err = session.ListService().DeleteListsByBoardID(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	err = session.CardService().DeleteCardsByBoardID(id)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handlePatchBoard handles requests to update a board.
func (h *BoardHandler) handlePatchBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req BoardUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	bu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	board, err := session.BoardService().UpdateBoard(id, bu)
	switch err {
	case nil:
		encodeJSON(w, board, http.StatusOK, h.Logger)
	case pulpe.ErrBoardNotFound:
		http.NotFound(w, r)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

// BoardCreateRequest is used to create a board.
type BoardCreateRequest struct {
	Name string `json:"name" valid:"required,stringlength(1|64)"`
}

// Validate board creation payload.
func (b *BoardCreateRequest) Validate() (*pulpe.BoardCreation, error) {
	b.Name = strings.TrimSpace(b.Name)
	err := validation.Validate(b)
	if err != nil {
		return nil, err
	}

	return &pulpe.BoardCreation{
		Name: b.Name,
	}, nil
}

// BoardUpdateRequest is used to update a board.
type BoardUpdateRequest struct {
	Name *string `json:"name" valid:"stringlength(1|64)"`
}

// Validate board update payload.
func (b *BoardUpdateRequest) Validate() (*pulpe.BoardUpdate, error) {
	if b.Name != nil {
		*b.Name = strings.TrimSpace(*b.Name)
	}

	err := validation.Validate(b)
	if b.Name != nil && *b.Name == "" {
		err = validation.AddError(err, "name", errors.New("name should not be empty"))
	}

	if err != nil {
		return nil, err
	}

	return &pulpe.BoardUpdate{
		Name: b.Name,
	}, nil
}

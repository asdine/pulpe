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

// registerBoardHandler register the BoardHandler routes.
func registerBoardHandler(router *httprouter.Router, c pulpeHttp.Connector) {
	h := boardHandler{
		connect: c,
		logger:  log.New(os.Stderr, "", log.LstdFlags),
	}

	router.HandlerFunc("GET", "/api/user/boards", h.handleGetBoards)
	router.HandlerFunc("POST", "/api/user/boards", h.handlePostBoard)
	router.GET("/api/boards/:owner/:board", h.handleGetBoard)
	router.DELETE("/api/boards/:id", h.handleDeleteBoard)
	router.PATCH("/api/boards/:id", h.handlePatchBoard)
}

// boardHandler represents an HTTP API handler for boards.
type boardHandler struct {
	connect pulpeHttp.Connector
	logger  *log.Logger
}

// handlePostBoard handles requests to create a new board.
func (h *boardHandler) handlePostBoard(w http.ResponseWriter, r *http.Request) {
	var req BoardCreateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	cr, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.connect(w, r)
	defer session.Close()

	board, err := session.BoardService().CreateBoard(cr)
	switch err {
	case nil:
		encodeJSON(w, board, http.StatusCreated, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handlePostBoard handles requests to create a new board.
func (h *boardHandler) handleGetBoards(w http.ResponseWriter, r *http.Request) {
	session := h.connect(w, r)
	defer session.Close()

	boards, err := session.BoardService().Boards()
	switch err {
	case nil:
		encodeJSON(w, boards, http.StatusOK, h.logger)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handleGetBoard handles requests to fetch a single board.
func (h *boardHandler) handleGetBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	owner := ps.ByName("owner")
	slug := ps.ByName("board")

	session := h.connect(w, r)
	defer session.Close()

	// Get the board and all of its lists and cards
	board, err := session.BoardService().BoardByOwnerAndSlug(owner, slug, pulpe.WithLists(), pulpe.WithCards())
	switch err {
	case nil:
		encodeJSON(w, board, http.StatusOK, h.logger)
	case pulpe.ErrBoardNotFound:
		http.NotFound(w, r)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handleDeleteBoard handles requests to delete a single board and all of its content.
func (h *boardHandler) handleDeleteBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	session := h.connect(w, r)
	defer session.Close()

	err := session.BoardService().DeleteBoard(id)
	switch err {
	case nil:
		w.WriteHeader(http.StatusNoContent)
	case pulpe.ErrBoardNotFound:
		http.NotFound(w, r)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
	}
}

// handlePatchBoard handles requests to update a board.
func (h *boardHandler) handlePatchBoard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var req BoardUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	bu, err := req.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.connect(w, r)
	defer session.Close()

	board, err := session.BoardService().UpdateBoard(id, bu)
	switch err {
	case nil:
		encodeJSON(w, board, http.StatusOK, h.logger)
	case pulpe.ErrBoardNotFound:
		http.NotFound(w, r)
	case pulpe.ErrUserAuthenticationFailed:
		Error(w, err, http.StatusUnauthorized, h.logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.logger)
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

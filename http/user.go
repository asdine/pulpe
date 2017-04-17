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

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(router *httprouter.Router, c pulpe.Client) *UserHandler {
	h := UserHandler{
		Router: router,
		Client: c,
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	h.POST("/register", h.handleUserRegistration)
	return &h
}

// UserHandler represents an HTTP API handler for users.
type UserHandler struct {
	*httprouter.Router

	Client pulpe.Client

	Logger *log.Logger
}

// handlePostUser handles requests to create a new user.
func (h *UserHandler) handleUserRegistration(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var payload UserRegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	ur, err := payload.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	session := h.Client.Connect()
	defer session.Close()

	user, err := session.UserService().CreateUser(ur)
	if err != nil {
		switch err {
		case pulpe.ErrUserEmailConflict:
			Error(w, validation.AddError(nil, "email", err), http.StatusBadRequest, h.Logger)
		default:
			Error(w, err, http.StatusInternalServerError, h.Logger)
		}
		return
	}

	encodeJSON(w, user, http.StatusCreated, h.Logger)
}

// UserRegistrationRequest is used to create a user.
type UserRegistrationRequest struct {
	FullName string `json:"name" valid:"required,stringlength(1|64)"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,stringlength(6|64)"`
}

// Validate user registration payload.
func (u *UserRegistrationRequest) Validate() (*pulpe.UserCreation, error) {
	u.FullName = strings.TrimSpace(u.FullName)
	err := validation.Validate(u)
	if err != nil {
		return nil, err
	}

	return &pulpe.UserCreation{
		FullName: u.FullName,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

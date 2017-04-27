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

// registerUserHandler register the userHandler routes.
func registerUserHandler(router *httprouter.Router, c *client) {
	h := userHandler{
		client: c,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}

	router.HandlerFunc("POST", "/api/register", h.handleUserRegistration)
	router.HandlerFunc("POST", "/api/login", h.handleUserLogin)
}

// userHandler represents an HTTP API handler for users.
type userHandler struct {
	client *client
	logger *log.Logger
}

// handleUserRegistration handles requests to create a new user.
func (h *userHandler) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var payload UserRegistrationRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	ur, err := payload.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	user, err := session.UserService().Register(ur)
	if err != nil {
		switch err {
		case pulpe.ErrUserEmailConflict:
			Error(w, validation.AddError(nil, "email", err), http.StatusBadRequest, h.logger)
		default:
			Error(w, err, http.StatusInternalServerError, h.logger)
		}
		return
	}

	us, err := session.UserSessionService().CreateSession(user)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.logger)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "pulpesid",
		Value:   us.ID,
		Expires: us.ExpiresAt,
		Path:    "/",
	})

	encodeJSON(w, user, http.StatusCreated, h.logger)
}

// handleUserLogin handles requests to login a user.
func (h *userHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var payload UserLoginRequest

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		Error(w, ErrInvalidJSON, http.StatusBadRequest, h.logger)
		return
	}

	err = payload.Validate()
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.logger)
		return
	}

	session := h.client.session(w, r)
	defer session.Close()

	us, err := session.UserSessionService().Login(payload.EmailOrLogin, payload.Password)
	if err != nil {
		switch err {
		case pulpe.ErrUserAuthenticationFailed:
			Error(w, err, http.StatusUnauthorized, h.logger)
		default:
			Error(w, err, http.StatusInternalServerError, h.logger)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "pulpesid",
		Value:   us.ID,
		Expires: us.ExpiresAt,
		Path:    "/",
	})

	w.WriteHeader(http.StatusCreated)
}

// UserRegistrationRequest is used to create a user.
type UserRegistrationRequest struct {
	FullName string `json:"fullName" valid:"required,stringlength(1|64)"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,stringlength(6|64)"`
}

// Validate user registration payload.
func (u *UserRegistrationRequest) Validate() (*pulpe.UserRegistration, error) {
	u.FullName = strings.TrimSpace(u.FullName)
	err := validation.Validate(u)
	if err != nil {
		return nil, err
	}

	return &pulpe.UserRegistration{
		FullName: u.FullName,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

// UserLoginRequest is used to login a user.
type UserLoginRequest struct {
	EmailOrLogin string `json:"login" valid:"required,stringlength(1|64)"`
	Password     string `json:"password" valid:"required,stringlength(1|64)"`
}

// Validate user login payload.
func (u *UserLoginRequest) Validate() error {
	return validation.Validate(u)
}

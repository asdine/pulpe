package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blankrobot/pulpe"
	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_Registration(t *testing.T) {
	t.Run("OK", testUserHandler_Registration_OK)
	t.Run("ErrInvalidJSON", testUserHandler_Registration_ErrInvalidJSON)
	t.Run("ErrValidation", testUserHandler_Registration_ErrValidation)
	t.Run("NotFound", testUserHandler_Registration_EmailConflict)
	t.Run("ErrInternal", testUserHandler_Registration_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testUserHandler_Registration_OK(t *testing.T) {
	c := mock.NewClient()

	c.UserService.CreateUserFn = func(u *pulpe.UserRegistration) (*pulpe.User, error) {
		require.Equal(t, &pulpe.UserRegistration{
			FullName: "Jon Snow",
			Email:    "jon.snow@wall.com",
			Password: "password",
		}, u)

		return &pulpe.User{
			ID:        "123",
			CreatedAt: mock.Now,
			Login:     "login",
			FullName:  "Jon Snow",
			Email:     "jon.snow@wall.com",
		}, nil
	}

	c.UserService.CreateSessionFn = func(u *pulpe.User) (*pulpe.UserSession, error) {
		require.Equal(t, &pulpe.User{
			ID:        "123",
			CreatedAt: mock.Now,
			Login:     "login",
			FullName:  "Jon Snow",
			Email:     "jon.snow@wall.com",
		}, u)

		return &pulpe.UserSession{
			ID:        "456",
			UpdatedAt: mock.Now,
			ExpiresAt: mock.Now.Add(10 * time.Minute),
		}, nil
	}

	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/signup", bytes.NewReader([]byte(`{
    "fullName": "Jon Snow",
    "email": "jon.snow@wall.com",
		"password": "password"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "pulpesid=456; Expires=Sat, 01 Jan 2000 00:10:00 GMT", w.HeaderMap.Get("Set-Cookie"))
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
		"fullName": "Jon Snow",
		"login": "login",
		"email": "jon.snow@wall.com",
		"createdAt": `+string(date)+`
	}`, w.Body.String())
}

func testUserHandler_Registration_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/signup", bytes.NewReader([]byte(`{
    "fullName": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testUserHandler_Registration_EmailConflict(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.UserService.CreateUserFn = func(user *pulpe.UserRegistration) (*pulpe.User, error) {
		return nil, pulpe.ErrUserEmailConflict
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/signup", bytes.NewReader([]byte(`{
		"fullName": "Jon Snow",
    "email": "jon.snow@wall.com",
		"password": "password"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func testUserHandler_Registration_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		c.UserService.CreateUserFn = func(user *pulpe.UserRegistration) (*pulpe.User, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/signup", bytes.NewReader([]byte(`{
			"fullName": "Jon Snow",
    "email": "jon.snow@wall.com",
		"password": "password"
		}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func testUserHandler_Registration_ErrValidation(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/signup", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login(t *testing.T) {
	t.Run("OK", testUserHandler_Login_OK)
	t.Run("ErrInvalidJSON", testUserHandler_Login_ErrInvalidJSON)
	t.Run("ErrValidation", testUserHandler_Login_ErrValidation)
	t.Run("NotFound", testUserHandler_Login_UserAuthenticationFailed)
	t.Run("ErrInternal", testUserHandler_Login_WithResponse(t, http.StatusInternalServerError, errors.New("unexpected error")))
}

func testUserHandler_Login_OK(t *testing.T) {
	c := mock.NewClient()

	c.UserService.AuthenticateFn = func(loginOrEmail, password string) (*pulpe.User, error) {
		require.Equal(t, "jonsnow", loginOrEmail)
		require.Equal(t, "password", password)

		return &pulpe.User{
			ID:        "123",
			CreatedAt: mock.Now,
			Login:     "login",
			FullName:  "Jon Snow",
			Email:     "jon.snow@wall.com",
		}, nil
	}

	c.UserService.CreateSessionFn = func(u *pulpe.User) (*pulpe.UserSession, error) {
		require.Equal(t, &pulpe.User{
			ID:        "123",
			CreatedAt: mock.Now,
			Login:     "login",
			FullName:  "Jon Snow",
			Email:     "jon.snow@wall.com",
		}, u)

		return &pulpe.UserSession{
			ID:        "456",
			UpdatedAt: mock.Now,
			ExpiresAt: mock.Now.Add(10 * time.Minute),
		}, nil
	}

	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader([]byte(`{
    "login": "jonsnow",
		"password": "password"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "pulpesid=456; Expires=Sat, 01 Jan 2000 00:10:00 GMT", w.HeaderMap.Get("Set-Cookie"))
	date, _ := mock.Now.MarshalJSON()
	require.JSONEq(t, `{
		"id": "123",
		"fullName": "Jon Snow",
		"login": "login",
		"email": "jon.snow@wall.com",
		"createdAt": `+string(date)+`
	}`, w.Body.String())
}

func testUserHandler_Login_ErrInvalidJSON(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader([]byte(`{
    "fullName": "12
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"err": "invalid_json"}`, w.Body.String())
}

func testUserHandler_Login_UserAuthenticationFailed(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	c.UserService.AuthenticateFn = func(loginOrEmail, password string) (*pulpe.User, error) {
		return nil, pulpe.ErrUserAuthenticationFailed
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader([]byte(`{
    "login": "jon.snow@wall.com",
		"password": "password"
  }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func testUserHandler_Login_WithResponse(t *testing.T, status int, err error) func(*testing.T) {
	return func(t *testing.T) {
		c := mock.NewClient()
		h := pulpeHttp.NewHandler(c)

		c.UserService.AuthenticateFn = func(loginOrEmail, password string) (*pulpe.User, error) {
			return nil, err
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader([]byte(`{
    "login": "jon.snow@wall.com",
		"password": "password"
		}`)))
		h.ServeHTTP(w, r)
		require.Equal(t, status, w.Code)
	}
}

func testUserHandler_Login_ErrValidation(t *testing.T) {
	c := mock.NewClient()
	h := pulpeHttp.NewHandler(c)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/v1/login", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

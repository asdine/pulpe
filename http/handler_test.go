package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/mock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	h := pulpeHttp.NewHandler(httprouter.New())
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/routethatdoesntexist", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCookieConnector(t *testing.T) {
	c := mock.NewClient()
	connect := pulpeHttp.NewCookieConnector(c)
	r, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	r.AddCookie(&http.Cookie{
		Name:  "pulpesid",
		Value: "token",
	})

	session := connect(nil, r)
	require.Equal(t, "token", session.(*mock.Session).AuthToken)
}

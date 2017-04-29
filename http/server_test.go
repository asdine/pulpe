package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/mock"
	"github.com/stretchr/testify/require"
)

func TestServeMux(t *testing.T) {
	mux := pulpeHttp.NewServeMux()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/routethatdoesntexist", bytes.NewReader([]byte(`{}`)))
	mux.ServeHTTP(w, r)
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

	session := connect(r)
	require.Equal(t, "token", session.(*mock.Session).AuthToken)
}

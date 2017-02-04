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

func TestHandler(t *testing.T) {
	h := pulpeHttp.NewHandler(mock.NewClient())
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/routethatdoesntexist", bytes.NewReader([]byte(`{}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

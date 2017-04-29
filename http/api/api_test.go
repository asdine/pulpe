package api_test

import (
	"net/http"

	"github.com/blankrobot/pulpe"
	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/http/api"
)

func newHandler(c pulpe.Client) http.Handler {
	mux := pulpeHttp.NewServeMux()
	connect := pulpeHttp.NewCookieConnector(c)
	api.Register(mux, connect)
	return mux
}

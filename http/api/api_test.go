package api_test

import (
	"net/http"

	"github.com/blankrobot/pulpe"
	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/http/api"
	"github.com/julienschmidt/httprouter"
)

func newHandler(c pulpe.Client) http.Handler {
	router := httprouter.New()
	connect := pulpeHttp.NewCookieConnector(c)
	api.Register(router, connect)
	return router
}

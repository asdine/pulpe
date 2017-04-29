package api

import (
	"encoding/json"
	"log"
	"net/http"

	pulpeHttp "github.com/blankrobot/pulpe/http"
	"github.com/julienschmidt/httprouter"
)

// Register all the routes and handlers to the given router
func Register(mux *pulpeHttp.ServeMux, connect pulpeHttp.Connector) {
	router := httprouter.New()
	registerBoardHandler(router, connect)
	registerCardHandler(router, connect)
	registerListHandler(router, connect)
	registerUserHandler(router, connect)

	mux.Handle("/api/", router)
}

// encodeJSON encodes v to w in JSON format. Error() is called if encoding fails.
func encodeJSON(w http.ResponseWriter, v interface{}, status int, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}

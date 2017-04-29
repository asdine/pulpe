package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/validation"
)

// HTTP errors
const (
	ErrInvalidJSON = pulpe.Error("invalid_json")
)

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error.
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = pulpe.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	switch {
	case validation.IsError(err):
		err = enc.Encode(&validationErrorResponse{
			Err:    "validation error",
			Fields: err,
		})
	default:
		err = enc.Encode(&errorResponse{Err: err.Error()})
	}

	if err != nil {
		logger.Println(err)
	}
}

// errorResponse is a generic response for sending an error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// validationErrorResponse is used for validation errors.
type validationErrorResponse struct {
	Err    string `json:"err,omitempty"`
	Fields error  `json:"fields"`
}

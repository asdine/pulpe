package pulpe

// General errors.
const (
	ErrInternal = Error("internal_error")
)

// Card errors
const (
	ErrCardNotFound = Error("card not found")
)

// List errors
const (
	ErrListNotFound = Error("list not found")
)

// Board errors
const (
	ErrBoardNotFound = Error("board not found")
)

// Error represents a Pulpe error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}

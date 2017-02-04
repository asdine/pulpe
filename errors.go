package pulpe

// General errors.
const (
	ErrInternal = Error("internal error")
)

// Card errors
const (
	ErrCardNotFound        = Error("card not found")
	ErrCardExists          = Error("card already exists")
	ErrCardIDRequired      = Error("card id required")
	ErrCardListIDRequired  = Error("card list id required")
	ErrCardBoardIDRequired = Error("card board id required")
)

// List errors
const (
	ErrListNotFound        = Error("list not found")
	ErrListExists          = Error("list already exists")
	ErrListIDRequired      = Error("list id required")
	ErrListBoardIDRequired = Error("list board id required")
)

// Board errors
const (
	ErrBoardNotFound   = Error("board not found")
	ErrBoardExists     = Error("board already exists")
	ErrBoardIDRequired = Error("board id required")
)

// Error represents a Pulpe error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}

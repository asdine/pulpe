package pulpe

// General errors.
const (
	ErrInternal = Error("internal_error")
)

// Error represents a Pulpe error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}

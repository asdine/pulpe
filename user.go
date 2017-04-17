package pulpe

import "time"

// User errors
const (
	ErrUserNotFound         = Error("user not found")
	ErrUserEmailConflict    = Error("email already exists")
	ErrAuthenticationFailed = Error("authentication failed")
)

// User informations.
type User struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	FullName  string     `json:"fullName"`
	Login     string     `json:"login"`
	Email     string     `json:"email"`
}

// UserCreation is used to create a User.
type UserCreation struct {
	FullName string
	Email    string
	Password string
}

// UserService represents a service for managing users.
type UserService interface {
	CreateUser(*UserCreation) (*User, error)
	User(id string) (*User, error)
	Authenticate(loginOrEmail, passwd string) (*User, error)
}

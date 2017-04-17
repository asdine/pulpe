package pulpe

import "time"

// User errors
const (
	ErrUserNotFound             = Error("user not found")
	ErrUserEmailConflict        = Error("email already exists")
	ErrUserAuthenticationFailed = Error("authentication failed")
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

// UserRegistration is used to register a User.
type UserRegistration struct {
	FullName string
	Email    string
	Password string
}

// UserSession is stored and represents a logged in user.
type UserSession struct {
	ID        string
	UserID    string
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// UserService represents a service for managing users.
type UserService interface {
	CreateUser(*UserRegistration) (*User, error)
	User(id string) (*User, error)
	Authenticate(loginOrEmail, passwd string) (*User, error)
	CreateSession(*User) (*UserSession, error)
}

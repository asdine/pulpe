package mock

import "github.com/blankrobot/pulpe"

// Ensure UserService implements pulpe.UserService.
var _ pulpe.UserService = new(UserService)

// UserService is a mock service that runs provided functions. Useful for testing.
type UserService struct {
	CreateUserFn      func(user *pulpe.UserRegistration) (*pulpe.User, error)
	CreateUserInvoked bool

	UserFn      func(id string) (*pulpe.User, error)
	UserInvoked bool

	AuthenticateFn      func(login, passwd string) (*pulpe.User, error)
	AuthenticateInvoked bool

	CreateSessionFn      func(user *pulpe.User) (*pulpe.UserSession, error)
	CreateSessionInvoked bool
}

// CreateUser runs CreateUserFn and sets CreateUserInvoked to true when invoked.
func (s *UserService) CreateUser(user *pulpe.UserRegistration) (*pulpe.User, error) {
	s.CreateUserInvoked = true
	return s.CreateUserFn(user)
}

// User runs UserFn and sets UserInvoked to true when invoked.
func (s *UserService) User(id string) (*pulpe.User, error) {
	s.UserInvoked = true
	return s.UserFn(id)
}

// Authenticate runs AuthenticateFn and sets AuthenticateInvoked to true when invoked.
func (s *UserService) Authenticate(login, passwd string) (*pulpe.User, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(login, passwd)
}

// CreateSession runs CreateSessionFn and sets CreateSessionInvoked to true when invoked.
func (s *UserService) CreateSession(user *pulpe.User) (*pulpe.UserSession, error) {
	s.CreateSessionInvoked = true
	return s.CreateSessionFn(user)
}

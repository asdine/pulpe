package mock

import "github.com/blankrobot/pulpe"

// Ensure UserService implements pulpe.UserService.
var _ pulpe.UserService = new(UserService)

// UserService is a mock service that runs provided functions. Useful for testing.
type UserService struct {
	CreateUserFn      func(user *pulpe.UserCreation) (*pulpe.User, error)
	CreateUserInvoked bool

	UserFn      func(id string) (*pulpe.User, error)
	UserInvoked bool
}

// CreateUser runs CreateUserFn and sets CreateUserInvoked to true when invoked.
func (s *UserService) CreateUser(user *pulpe.UserCreation) (*pulpe.User, error) {
	s.CreateUserInvoked = true
	return s.CreateUserFn(user)
}

// User runs UserFn and sets UserInvoked to true when invoked.
func (s *UserService) User(id string) (*pulpe.User, error) {
	s.UserInvoked = true
	return s.UserFn(id)
}

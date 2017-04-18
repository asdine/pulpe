package mock

import "github.com/blankrobot/pulpe"

// Ensure UserService implements pulpe.UserService.
var _ pulpe.UserService = new(UserService)

// UserService is a mock service that runs provided functions. Useful for testing.
type UserService struct {
	RegisterFn      func(user *pulpe.UserRegistration) (*pulpe.User, error)
	RegisterInvoked bool

	UserFn      func(id string) (*pulpe.User, error)
	UserInvoked bool

	LoginFn      func(login, passwd string) (*pulpe.User, error)
	LoginInvoked bool
}

// Register runs RegisterFn and sets RegisterInvoked to true when invoked.
func (s *UserService) Register(user *pulpe.UserRegistration) (*pulpe.User, error) {
	s.RegisterInvoked = true
	return s.RegisterFn(user)
}

// User runs UserFn and sets UserInvoked to true when invoked.
func (s *UserService) User(id string) (*pulpe.User, error) {
	s.UserInvoked = true
	return s.UserFn(id)
}

// Login runs LoginFn and sets LoginInvoked to true when invoked.
func (s *UserService) Login(login, passwd string) (*pulpe.User, error) {
	s.LoginInvoked = true
	return s.LoginFn(login, passwd)
}

// Ensure UserSessionService implements pulpe.UserSessionService.
var _ pulpe.UserSessionService = new(UserSessionService)

// UserSessionService is a mock service that runs provided functions. Useful for testing.
type UserSessionService struct {
	CreateSessionFn      func(user *pulpe.User) (*pulpe.UserSession, error)
	CreateSessionInvoked bool

	GetSessionFn      func(sid string) (*pulpe.UserSession, error)
	GetSessionInvoked bool
}

// CreateSession runs CreateSessionFn and sets CreateSessionInvoked to true when invoked.
func (s *UserSessionService) CreateSession(user *pulpe.User) (*pulpe.UserSession, error) {
	s.CreateSessionInvoked = true
	return s.CreateSessionFn(user)
}

// GetSession runs GetSessionFn and sets GetSessionInvoked to true when invoked.
func (s *UserSessionService) GetSession(sid string) (*pulpe.UserSession, error) {
	s.GetSessionInvoked = true
	return s.GetSessionFn(sid)
}

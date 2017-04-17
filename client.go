package pulpe

// Client creates a connection to the services.
type Client interface {
	Connect() Session
}

// Session represents a connection to the services.
type Session interface {
	CardService() CardService
	ListService() ListService
	BoardService() BoardService
	UserService() UserService
	Close() error
}

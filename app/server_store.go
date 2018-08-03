package app

// ServerStore is the interface that covers the management
// of the persistence layer for the servers
type ServerStore interface {
	Save(*Server) error
	RetrieveServerByName(string) (*Server, error)
}

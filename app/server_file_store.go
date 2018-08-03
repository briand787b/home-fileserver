package app

import (
	"os"

	"github.com/pkg/errors"
)

const (
	defaultServersPath = "servers.json"
)

// serverFileStore implements the ServerStore
// interface with flat file storage in JSON
type serverFileStore struct {
	file *os.File
	Servers
}

// NewServerFileStore sets up and returns a new
// ServerStore using flat file storage
func NewServerFileStore(path string) (ServerStore, error) {
	if path == "" {
		path = defaultServersPath
	}

	fd, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0600,
	)

	if err != nil {
		return errors.Wrap(err,
			"could not open file for file store",
		)
	}
}

// Save saves a Server to persistent storage
func (sfs *serverFileStore) Save(s *Server) error {
	return nil
}

func (sfs *serverFileStore) GetAllServers() (Servers, error) {
	// if s.Servers {

	// }

	return s.Servers
}

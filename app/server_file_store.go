package app

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

const (
	serversPath = "servers.json"
)

// fileServerStore implements the ServerStore
// interface with flat file storage in JSON
//
// The file which persists data will only
// be opened when actively reading or writing.
// Since reads/writes only happen when things
// change, and the server configs should change
// infrequently, every time that file is opened
// it will be closed within the same invocation
// in which it is opened
type fileServerStore map[string]*Server

// NewFileServerStore sets up and returns a new
// ServerStore using flat file storage
func NewFileServerStore() (ServerStore, error) {
	return newFileServerStore(serversPath)
}

// newFileServerStore exists to facilitate testing of the
// fileServerStore struct
func newFileServerStore(path string) (ServerStore, error) {
	fss := &fileServerStore{}
	return fss, fss.getState(path)
}

// saveState saves the current state of fss to disk
func (fss *fileServerStore) saveState(filename string) error {
	fd, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err,
			"could not open file for file store",
		)
	}

	// do i need to truncate here?

	if err := json.NewEncoder(fd).Encode(fss); err != nil {
		return errors.Wrap(err,
			"could not encode fss into file as JSON",
		)
	}

	return nil
}

func (fss *fileServerStore) getState(path string) error {
	fd, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_RDWR,
		0600,
	)
	if err != nil {
		return errors.Wrap(err,
			"could not open file for file store")
	}
	defer fd.Close()

	fStat, err := fd.Stat()
	if err != nil {
		return errors.Wrapf(err,
			"could not get stats of file %s", path)
	}

	if fStat.Size() == 0 {
		// file is empty, fss should be empty too
		fss = &fileServerStore{}
		return nil
	}

	if err := json.NewDecoder(fd).Decode(&fss); err != nil {
		return errors.Wrap(err,
			"could not decode json in file to fileServerStore")
	}

	return nil
}

// Save saves a Server to persistent storage
func (fss *fileServerStore) Save(s *Server) error {
	return fss.save(s)
}

func (fss *fileServerStore) save(s *Server) error {
	return nil
}

func (fss *fileServerStore) GetAllServers() ([]Server, error) {
	// if s.Servers {

	// }

	return nil, nil
}

func (fss *fileServerStore) GetServerByName(name string) (*Server, error) {
	return nil, nil
}

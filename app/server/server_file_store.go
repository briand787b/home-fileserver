package server

import (
	"encoding/json"
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	serversPath = "servers.json"
)

var (
	// ErrServerNotFound is the type of error returned when the
	// ServerStore cannot find the named server
	ErrServerNotFound = errors.New("Server does not exist in store")
)

// fileStore implements the Store
// interface with flat file storage in JSON
//
// The file which persists data will only
// be opened when actively reading or writing.
// Since reads/writes only happen when things
// change, and the server configs should change
// infrequently, every time that file is opened
// it will be closed within the same invocation
// in which it is opened
type fileStore map[string]*Server

// NewFileStore sets up and returns a new
// Store using flat file storage
func NewFileStore() (Store, error) {
	return newFileStore(serversPath)
}

// newFileStore exists to facilitate testing of the
// fileStore struct
func newFileStore(path string) (fileStore, error) {
	fss := fileStore{}
	return fss, fss.getState(path)
}

// saveState adds the current state of fss to disk
func (fss fileStore) saveState(filename string) error {
	fd, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return errors.Wrap(err,
			"could not open file for file store",
		)
	}
	defer fd.Close()

	if err := fd.Truncate(0); err != nil {
		return errors.Wrap(err,
			"cannot truncate server store file")
	}

	if err := json.NewEncoder(fd).Encode(&fss); err != nil {
		return errors.Wrap(err,
			"could not encode fss into file as JSON",
		)
	}

	return nil
}

func (fss fileStore) getState(path string) error {
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
		fss = fileStore{}
		return nil
	}

	if err := json.NewDecoder(fd).Decode(&fss); err != nil {
		return errors.Wrap(err,
			"could not decode json in file to fileStore")
	}

	return nil
}

// Add adds a server and saves it to persistent storage
func (fss fileStore) AddServer(s *Server) error {
	return fss.add(s, serversPath)
}

func (fss fileStore) add(s *Server, path string) error {
	if s.Name == "" {
		return errors.New("Server must have name")
	}

	if net.ParseIP(s.IP) == nil {
		return errors.New("IP address is not valid")
	}

	if s.WorkingDir == "" {
		return errors.New("WorkingDir cannot be empty")
	}

	fss[s.Name] = s
	return fss.saveState(path)
}

func (fss fileStore) RemoveServer(name string) error {
	return fss.removeServer(name, serversPath)
}

func (fss fileStore) removeServer(name, path string) error {
	delete(fss, name)
	return fss.saveState(path)
}

func (fss fileStore) GetAllServers() (ss []*Server) {
	for _, s := range fss {
		ss = append(ss, s)
	}

	return
}

func (fss fileStore) GetServerByName(name string) (*Server, error) {
	s, ok := fss[name]
	if ok {
		return s, nil
	}

	return nil, ErrServerNotFound
}

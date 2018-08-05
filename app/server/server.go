package server

// Server represents a physical server that hosts
// and/or consumes content.
//
// FileList is not persisted because it is dynamic
// by the nature of this program and its accuracy
// can only be guaranteed on fresh retrievals from
// the server in question
type Server struct {
	Name       string   `json:"name"`
	IP         string   `json:"ip_addr"`
	WorkingDir string   `json:"working_dir"` // absolute path
	FileList   []string `json:"file_list"`
	Local      bool     `json:"local"`
}

// NewServer instantiates a new server.
// It does NOT save it to the persistent
// storage
func NewServer(name, ip string) *Server {
	return &Server{
		Name: name,
		IP:   ip,
	}
}

func NewLocalServer(dir string) *Server {
	// get files in dir

	// get host ip

	// get name from host

	// set local to true

	return nil
}

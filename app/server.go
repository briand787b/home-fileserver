package app

// Server represents a physical server that hosts
// and/or consumes content.
type Server struct {
	Name       string
	IP         string
	WorkingDir string // absolute path
	FileList   []string
}

// NewServer instantiates a new server.
// It does NOT save it to the persistent
// storage
func NewServer(name, ip, dir string) *Server {
	return &Server{
		IP:         ip,
		WorkingDir: dir,
	}
}

// Servers represents the set of known servers
// on which the hfs binary is installed
type Servers map[string]*Server

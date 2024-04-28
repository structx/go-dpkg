package domain

// Logger configuration
type Logger struct {
	Level    string `hcl:"log_level"`
	Path     string `hcl:"log_path"`
	RaftPath string `hcl:"raft_log_path"`
}

// Raft configuration
type Raft struct {
	Bootstrap bool   `hcl:"bootstrap"`
	LocalID   string `hcl:"local_id"`
	BaseDir   string `hcl:"base_dir"`
}

// Ports configuration
type Ports struct {
	HTTP int `hcl:"http"`
	GRPC int `hcl:"grpc"`
}

// Server configuration
type Server struct {
	BindAddr       string `hcl:"bind_addr"`
	Ports          Ports  `hcl:"ports,block"`
	DefaultTimeout int64  `hcl:"default_timeout"`
}

// Chain configuration
type Chain struct {
	BaseDir string `hcl:"base_dir"`
}

// Messenger message broker configuration
type Messenger struct {
	ServerAddr string `hcl:"server_addr"`
}

// Config service configuration interface
type Config interface {
	GetServer() Server
	GetRaft() Raft
	GetLogger() Logger
	GetChain() Chain
	GetMessenger() Messenger
}

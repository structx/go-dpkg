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

// AccessControl configuration
type AccessControl struct {
	ServerAddr string `hcl:"server_addr"`
}

// DistributedHashTable configuration
type DistributedHashTable struct {
	BindAddr      string  `hcl:"bind_addr"`
	AdvertiseAddr *string `hcl:"advertise_addr"`
	Ports         struct {
		GRPC int `hcl:"grpc"`
	} `hcl:"ports"`
}

// Config service configuration interface
//
//go:generate mockery --name Config
type Config interface {
	// GetServer getter server configuration
	GetServer() Server
	// GetRaft getter raft configuration
	GetRaft() *Raft
	// GetLogger getter logger configuration
	GetLogger() Logger
	// GetChain getter chain configuration
	GetChain() *Chain
	// GetMessenger getter messenger configuration
	GetMessenger() *Messenger
	// GetAccessControl getter access control configuration
	GetAccessControl() *AccessControl
	// GetDistributedHashTable getter distributed hash table configuration
	GetDistributedHashTable() *DistributedHashTable
}

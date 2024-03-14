package kv

// Config kv database configuration
type Config struct {
	Dir string `env:"KV_DIR"`
}

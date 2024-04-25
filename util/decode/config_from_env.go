// Package decode files into hcl
package decode

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	hcl2 "github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/trevatk/go-pkg/domain"
)

// ConfigFromEnv decode service configuration from config file
func ConfigFromEnv(cfg domain.Config) error {

	configFile := os.Getenv("DSERVICE_CONFIG")
	if configFile == "" {
		return errors.New("$DSERVICE_CONFIG must be set")
	}

	if err := hcl2.DecodeFile(
		filepath.Clean(configFile),
		nil,
		cfg,
	); err != nil {
		return fmt.Errorf("failed decode config file %v", err)
	}

	return nil
}

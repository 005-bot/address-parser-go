package address

import (
	"fmt"
	"os"
	"path/filepath"
)

func createTempDir() (string, error) {
	dir, err := os.MkdirTemp("", "address-parser-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}

	return dir, nil
}

func writeEmbeddedDB(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, streetsDB, 0600); err != nil {
		return fmt.Errorf("write db file: %w", err)
	}

	return nil
}

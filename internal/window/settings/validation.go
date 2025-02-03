package settings

import (
	"fmt"
	"os"
)

func requirePathIsFolder(s string) error {
	if s == "" {
		return nil
	}

	info, err := os.Stat(s)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", s)
	}

	return nil
}

func requirePathIsFile(s string) error {
	if s == "" {
		return nil
	}

	info, err := os.Stat(s)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", s)
	}

	return nil
}

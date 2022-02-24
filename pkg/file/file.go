package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func ForToken() (*os.File, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("getting user config dir: %w", err)
	}

	file, err := os.Create(filepath.Join(configDir, "kainotomia", "token"))
	if err != nil {
		return nil, fmt.Errorf("creating file for token: %w", err)
	}
	return file, nil
}

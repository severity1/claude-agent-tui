package tui_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryStructure(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	requiredDirs := []string{
		"example/chat",
		"component/output/streamtext",
		"component/input/chatinput",
		"component/input/permission",
		"adapter",
		"layout/chat",
	}

	for _, dir := range requiredDirs {
		t.Run(dir, func(t *testing.T) {
			path := filepath.Join(root, dir)
			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				t.Errorf("required directory does not exist: %s", dir)
				return
			}
			if err != nil {
				t.Errorf("error checking directory %s: %v", dir, err)
				return
			}
			if !info.IsDir() {
				t.Errorf("path exists but is not a directory: %s", dir)
			}
		})
	}
}

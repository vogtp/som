package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/log"
)

// IsDockerContainer checks if running in a docker env
func IsDockerContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

// EnsureOutFolder creates a folder if it does not exist
func EnsureOutFolder(outFolder string) error {
	if outFolder != "" {
		if _, err := os.Stat(outFolder); os.IsNotExist(err) {
			err := os.MkdirAll(outFolder, os.ModePerm)
			if err != nil {
				hcl.Error("Cannot create output folder", "folder", outFolder, log.Error, err)
				return fmt.Errorf("Cannot create output folder %q: %v", outFolder, err)
			}
		}
	}
	if outFolder != "" && !strings.HasSuffix(outFolder, "/") {
		outFolder = outFolder + "/"
	}
	fn := outFolder + "test.txt"
	if err := os.WriteFile(fn, []byte("Testoutput"), 0644); err != nil {
		dir, err2 := os.Getwd()
		if err2 != nil {
			hcl.Error("Getwd failed", log.Error, err)
		}
		hcl.Error("Cannot wirte", log.Error, err, "dir", dir)
		return fmt.Errorf("Cannot wirte out: %v", err)
	}
	return os.Remove(fn)
}

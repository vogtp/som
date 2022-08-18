package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vogtp/go-hcl"
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
				hcl.Errorf("Cannot create output folder %q: %v", outFolder, err)
				return fmt.Errorf("Cannot create output folder %q: %v", outFolder, err)
			}
		}
	}
	if outFolder != "" && !strings.HasSuffix(outFolder, "/") {
		outFolder = outFolder + "/"
	}
	fn := outFolder + "test.txt"
	if err := ioutil.WriteFile(fn, []byte("Testoutput"), 0644); err != nil {
		dir, err2 := os.Getwd()
		if err2 != nil {
			hcl.Errorf("Getwd: %v", err)
		}
		hcl.Warnf("cwd: %v", dir)
		hcl.Errorf("Cannot wirte out: %v", err)
		return fmt.Errorf("Cannot wirte out: %v", err)
	}
	return os.Remove(fn)
}

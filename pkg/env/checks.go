package env

import (
	"os"
	"strings"
)

// IsGoRun checks if run by go run
// it does this by checking arg[0]
func IsGoRun() bool {
	if IsGoTest() {
		return false
	}
	a := os.Args[0]
	if !strings.Contains(a, "go-build") {
		return false
	}
	i1 := strings.Index(a, "/go-build")
	i2 := strings.Index(a, "\\go-build")
	if i1+i2 < 0 {
		return false
	}
	return true
	// s := string(a[i+len("go-build")])
	// _, err := strconv.Atoi(s)
	// slog.Info("go build", "s", s, "a", a, "err", err)
	// return err == nil
}

// IsGoTest checks if run by go test
// it does this by checking arg[0]
func IsGoTest() bool {
	a := os.Args[0]
	if strings.HasSuffix(a, ".test") {
		return true
	}
	if strings.HasSuffix(a, ".test.exe") {
		return true
	}
	if strings.HasSuffix(a, "__debug_bin") {
		return true
	}
	return false
}

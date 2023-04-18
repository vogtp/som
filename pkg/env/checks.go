package env

import (
	"os"
	"strconv"
	"strings"
)

// IsGoRun checks if run by go run
// it does this by checking arg[0]
func IsGoRun() bool {
	if IsGoTest() {
		return false
	}
	a := os.Args[0]
	i := strings.Index(a, "go-build")
	if i == -1 {
		return false
	}
	i1 := strings.Index(a, "/go-build")
	i2 := strings.Index(a, "\\go-build")
	if i1+i2 < 0 {
		return false
	}
	s := string(a[i+len("go-build")])
	_, err := strconv.Atoi(s)
	return err == nil
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

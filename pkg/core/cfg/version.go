package cfg

import "fmt"

const (
	// VersionMajor major version
	VersionMajor = 0
	// VersionMinor minor version
	VersionMinor = 10
	// VersionPatch patch level
	VersionPatch = 1
)

var (
	// BuildInfo contains the build timestamp
	BuildInfo = "development"
	// Version string
	Version = fmt.Sprintf("%v.%v.%v (%v)", VersionMajor, VersionMinor, VersionPatch, BuildInfo)
)

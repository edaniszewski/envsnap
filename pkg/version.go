package pkg

import "runtime"

// Build-time arguments are used to populate these variables to generate
// build-specific version information.
var (
	Version   string
	Commit    string
	Tag       string
	GoVersion string
	BuildDate string
)

// versionInfo contains runtime and build time information about the application.
type versionInfo struct {
	Version   string
	Commit    string
	Tag       string
	GoVersion string
	BuildDate string
	Compiler  string
	OS        string
	Arch      string
}

// newVersionInfo creates a struct which holds all of the runtime and
// build-time supplied variables describing the version and build state
// for the application.
func newVersionInfo() versionInfo {
	return versionInfo{
		Version:   Version,
		Commit:    Commit,
		Tag:       Tag,
		GoVersion: GoVersion,
		BuildDate: BuildDate,
		Compiler:  runtime.Compiler,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

package pkg

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionInfo(t *testing.T) {
	// The versionInfo should not have any of the build time parameters
	// set for this test.
	vi := newVersionInfo()

	assert.Equal(t, "", vi.Version)
	assert.Equal(t, "", vi.Commit)
	assert.Equal(t, "", vi.Tag)
	assert.Equal(t, "", vi.GoVersion)
	assert.Equal(t, "", vi.BuildDate)
	assert.Equal(t, runtime.Compiler, vi.Compiler)
	assert.Equal(t, runtime.GOOS, vi.OS)
	assert.Equal(t, runtime.GOARCH, vi.Arch)
}

package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemConfig_Render_NoOpts(t *testing.T) {
	cfg := SystemConfig{}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, SystemResult{}, r)

	res := r.(SystemResult)
	assert.Empty(t, res.OS)
	assert.Empty(t, res.Arch)
	assert.Empty(t, res.CPUs)
	assert.Empty(t, res.Kernel)
	assert.Empty(t, res.KernelVersion)
	assert.Empty(t, res.Processor)
}

func TestSystemConfig_Render_BadOpt(t *testing.T) {
	cfg := SystemConfig{
		Core: []string{"not-an-option"},
	}

	r, err := cfg.Render()
	assert.Error(t, err)
	assert.IsType(t, SystemResult{}, r)
}

func TestSystemConfig_Render(t *testing.T) {
	cfg := SystemConfig{
		Core: []string{
			"os", "arch", "cpus", "kernel", "kernel_version", "processor",
		},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, SystemResult{}, r)

	res := r.(SystemResult)
	assert.NotEmpty(t, res.OS)
	assert.NotEmpty(t, res.Arch)
	assert.NotEmpty(t, res.CPUs)
	assert.NotEmpty(t, res.Kernel)
	assert.NotEmpty(t, res.KernelVersion)
	assert.NotEmpty(t, res.Processor)
}

func TestSystemConfig_Render_Alternatives(t *testing.T) {
	cfg := SystemConfig{
		Core: []string{
			"kernel-version",
		},
	}

	r, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, SystemResult{}, r)

	res := r.(SystemResult)
	assert.Empty(t, res.OS)
	assert.Empty(t, res.Arch)
	assert.Empty(t, res.CPUs)
	assert.Empty(t, res.Kernel)
	assert.NotEmpty(t, res.KernelVersion)
	assert.Empty(t, res.Processor)
}

func TestNewSystemResult(t *testing.T) {
	r := NewSystemResult()
	assert.Empty(t, r.OS)
	assert.Empty(t, r.Arch)
	assert.Empty(t, r.Kernel)
	assert.Empty(t, r.KernelVersion)
	assert.Empty(t, r.Processor)
	assert.Empty(t, r.CPUs)
}

func TestSystemResult_IsEmpty(t *testing.T) {
	r := NewSystemResult()
	assert.True(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty2(t *testing.T) {
	r := NewSystemResult()
	r.CPUs = 1

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty3(t *testing.T) {
	r := NewSystemResult()
	r.OS = "os"

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty4(t *testing.T) {
	r := NewSystemResult()
	r.Arch = "arch"

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty5(t *testing.T) {
	r := NewSystemResult()
	r.Kernel = "kernel"

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty6(t *testing.T) {
	r := NewSystemResult()
	r.KernelVersion = "ver"

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_IsEmpty7(t *testing.T) {
	r := NewSystemResult()
	r.Processor = "proc"

	assert.False(t, r.IsEmpty())
}

func TestSystemResult_Markdown(t *testing.T) {
	r := NewSystemResult()
	r.OS = "darwin"
	r.Arch = "x86_64"
	r.CPUs = 12
	r.Kernel = "Darwin"
	r.KernelVersion = "19.0.0"
	r.Processor = "i386"

	data, err := r.Markdown()
	assert.NoError(t, err)

	expected := "**System**\n- _os_: darwin\n- _arch_: x86_64\n- _cpus_: 12\n- _kernel_: Darwin\n- _kernel version_: 19.0.0\n- _processor_: i386\n"
	assert.Equal(t, expected, string(data))
}

func TestSystemResult_Markdown_Empty(t *testing.T) {
	r := NewSystemResult()

	data, err := r.Markdown()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestSystemResult_Plaintext(t *testing.T) {
	r := NewSystemResult()
	r.OS = "darwin"
	r.Arch = "x86_64"
	r.CPUs = 12
	r.Kernel = "Darwin"
	r.KernelVersion = "19.0.0"
	r.Processor = "i386"

	data, err := r.Plaintext()
	assert.NoError(t, err)

	expected := "System\n------\nos:             darwin\narch:           x86_64\ncpus:           12\nkernel:         Darwin\nkernel version: 19.0.0\nprocessor:      i386\n"
	assert.Equal(t, expected, string(data))
}

func TestSystemResult_Plaintext_Empty(t *testing.T) {
	r := NewSystemResult()

	data, err := r.Plaintext()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestSystemResult_JSON(t *testing.T) {
	r := NewSystemResult()
	r.OS = "darwin"
	r.Arch = "x86_64"
	r.CPUs = 12
	r.Kernel = "Darwin"
	r.KernelVersion = "19.0.0"
	r.Processor = "i386"

	data, err := r.JSON()
	assert.NoError(t, err)

	expected := `{"os":"darwin","arch":"x86_64","cpus":12,"kernel":"Darwin","kernel_version":"19.0.0","processor":"i386"}`
	assert.Equal(t, expected, string(data))
}

func TestSystemResult_JSON_Empty(t *testing.T) {
	r := NewSystemResult()

	data, err := r.JSON()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestSystemResult_YAML(t *testing.T) {
	r := NewSystemResult()
	r.OS = "darwin"
	r.Arch = "x86_64"
	r.CPUs = 12
	r.Kernel = "Darwin"
	r.KernelVersion = "19.0.0"
	r.Processor = "i386"

	data, err := r.YAML()
	assert.NoError(t, err)

	expected := "os: darwin\narch: x86_64\ncpus: 12\nkernel: Darwin\nkernel_version: 19.0.0\nprocessor: i386\n"
	assert.Equal(t, expected, string(data))
}

func TestSystemResult_YAML_Empty(t *testing.T) {
	r := NewSystemResult()

	data, err := r.YAML()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

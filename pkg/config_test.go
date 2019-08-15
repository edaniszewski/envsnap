package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV1EnvsnapConfig_All(t *testing.T) {
	cfg := V1EnvsnapConfig{}

	all := cfg.All()
	assert.Len(t, all, 5)
	assert.IsType(t, SystemConfig{}, all[0])
	assert.IsType(t, EnvConfig{}, all[1])
	assert.IsType(t, ExecConfig{}, all[2])
	assert.IsType(t, PythonConfig{}, all[3])
	assert.IsType(t, GolangConfig{}, all[4])
}

func TestV1EnvsnapConfig_Render_Err(t *testing.T) {
	cfg := V1EnvsnapConfig{
		System: SystemConfig{
			Core: []string{"foobar"},
		},
	}

	out, err := cfg.Render()
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestV1EnvsnapConfig_Render_Err2(t *testing.T) {
	cfg := V1EnvsnapConfig{
		Python: PythonConfig{
			Core: []string{"foobar"},
		},
	}

	out, err := cfg.Render()
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestV1EnvsnapConfig_Render_Err3(t *testing.T) {
	cfg := V1EnvsnapConfig{
		Golang: GolangConfig{
			Core: []string{"foobar"},
		},
	}

	out, err := cfg.Render()
	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestV1EnvsnapConfig_Render_Empty(t *testing.T) {
	cfg := V1EnvsnapConfig{}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, &V1EnvsnapResult{}, out)

	res := out.(*V1EnvsnapResult)
	assert.True(t, res.System.IsEmpty())
	assert.True(t, res.Environment.IsEmpty())
	assert.True(t, res.Exec.IsEmpty())
	assert.True(t, res.Python.IsEmpty())
	assert.True(t, res.Golang.IsEmpty())
}

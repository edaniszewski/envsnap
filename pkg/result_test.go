package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewV1EnvsnapResult(t *testing.T) {
	v1 := NewV1EnvsnapResult()
	assert.Equal(t, os.Stdout, v1.out)
	assert.Nil(t, v1.Environment)
	assert.Nil(t, v1.Exec)
	assert.Nil(t, v1.Golang)
	assert.Nil(t, v1.Python)
	assert.Nil(t, v1.System)
}

func TestV1EnvsnapResult_Results(t *testing.T) {
	v1 := NewV1EnvsnapResult()
	v1.System = NewSystemResult()
	v1.Environment = NewEnvResult()
	v1.Exec = NewExecResult()
	v1.Python = NewPythonResult()
	v1.Golang = NewGolangResult()

	res := v1.Results()
	assert.Len(t, res, 5)
	assert.IsType(t, SystemResult{}, res[0])
	assert.IsType(t, EnvResult{}, res[1])
	assert.IsType(t, ExecResult{}, res[2])
	assert.IsType(t, PythonResult{}, res[3])
	assert.IsType(t, GolangResult{}, res[4])
}

func TestV1EnvsnapResult_String_Markdown(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	data, err := v1.String("markdown")
	assert.NoError(t, err)

	assert.Equal(t, "#### Environment\n\n**System**\n- _os_: testOS\n", data)
}

func TestV1EnvsnapResult_String_Plaintext(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	data, err := v1.String("plaintext")
	assert.NoError(t, err)

	assert.Equal(t, "System\n------\nos:             testOS\n", data)
}

func TestV1EnvsnapResult_String_JSON(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	data, err := v1.String("json")
	assert.NoError(t, err)

	assert.Equal(t, `{"system":{"os":"testOS"}}`, data)
}

func TestV1EnvsnapResult_String_YAML(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	data, err := v1.String("yaml")
	assert.NoError(t, err)

	assert.Equal(t, "system:\n  os: testOS\n", data)
}

func TestV1EnvsnapResult_String_UnsupportedFmt(t *testing.T) {
	v1 := NewV1EnvsnapResult()

	str, err := v1.String("foobar")
	assert.Empty(t, str)
	assert.Error(t, err)
	assert.Equal(t, ErrUnsupportedFormat, err)
}

func TestV1EnvsnapResult_Write(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	file, err := ioutil.TempFile("", "envsnap-test")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	err = v1.Write(file.Name(), "json")
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(file.Name())
	assert.NoError(t, err)
	assert.Equal(t, `{"system":{"os":"testOS"}}`, string(data))
}

func TestV1EnvsnapResult_Print(t *testing.T) {
	sys := NewSystemResult()
	sys.OS = "testOS"

	v1 := NewV1EnvsnapResult()
	v1.System = sys

	out := bytes.Buffer{}
	v1.out = &out

	err := v1.Print("json")
	assert.NoError(t, err)

	assert.Equal(t, "{\"system\":{\"os\":\"testOS\"}}\n", out.String())
}

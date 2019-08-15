package pkg

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestPythonConfig_Render(t *testing.T) {
	defer cliWarnings.Clear()

	// Since `python` may not be installed on the machine that is running
	// the tests, check to see whether it exists, which determines the outcome
	// of the test.
	hasBin := binExists("python")

	cfg := PythonConfig{
		Core: []string{"version"},
	}

	out, err := cfg.Render()
	if hasBin {
		assert.NoError(t, err)
		assert.IsType(t, PythonResult{}, out)
		res := out.(PythonResult)

		if res.Version == "" {
			assert.Contains(t, cliWarnings.Warnings, "python.core.version")
			assert.Len(t, cliWarnings.Warnings["python.core.version"], 1)
		}
		assert.Empty(t, res.VersionPy2)
		assert.Empty(t, res.VersionPy3)
		assert.Empty(t, res.Deps)

	} else {
		assert.NoError(t, err)
		assert.Contains(t, cliWarnings.Warnings, "python.core.version")
		assert.Len(t, cliWarnings.Warnings["python.core.version"], 1)
	}
}

func TestPythonConfig_Render2(t *testing.T) {
	defer cliWarnings.Clear()

	// Since `python2` may not be installed on the machine that is running
	// the tests, check to see whether it exists, which determines the outcome
	// of the test.
	hasBin := binExists("python2")

	cfg := PythonConfig{
		Core: []string{"py2"},
	}

	out, err := cfg.Render()
	if hasBin {
		assert.NoError(t, err)
		assert.IsType(t, PythonResult{}, out)
		res := out.(PythonResult)

		if res.VersionPy2 == "" {
			assert.Contains(t, cliWarnings.Warnings, "python.core.py2")
			assert.Len(t, cliWarnings.Warnings["python.core.py2"], 1)
		}
		assert.Empty(t, res.Version)
		assert.Empty(t, res.VersionPy3)
		assert.Empty(t, res.Deps)

	} else {
		assert.NoError(t, err)
		assert.Contains(t, cliWarnings.Warnings, "python.core.py2")
		assert.Len(t, cliWarnings.Warnings["python.core.py2"], 1)
	}
}

func TestPythonConfig_Render3(t *testing.T) {
	defer cliWarnings.Clear()

	// Since `python3` may not be installed on the machine that is running
	// the tests, check to see whether it exists, which determines the outcome
	// of the test.
	hasBin := binExists("python3")

	cfg := PythonConfig{
		Core: []string{"py3"},
	}

	out, err := cfg.Render()
	if hasBin {
		assert.NoError(t, err)
		assert.IsType(t, PythonResult{}, out)
		res := out.(PythonResult)

		if res.VersionPy3 == "" {
			assert.Contains(t, cliWarnings.Warnings, "python.core.py3")
			assert.Len(t, cliWarnings.Warnings["python.core.py3"], 1)
		}
		assert.Empty(t, res.Version)
		assert.Empty(t, res.VersionPy2)
		assert.Empty(t, res.Deps)

	} else {
		assert.NoError(t, err)
		assert.Contains(t, cliWarnings.Warnings, "python.core.py3")
		assert.Len(t, cliWarnings.Warnings["python.core.py3"], 1)
	}
}

func TestPythonConfig_Render4(t *testing.T) {
	defer cliWarnings.Clear()

	// Since `pip` may not be installed on the machine that is running
	// the tests, check to see whether it exists, which determines the outcome
	// of the test.
	hasBin := binExists("pip")

	cfg := PythonConfig{
		Deps: DependenciesConfig{
			Packages: []string{"setuptools"},
		},
	}

	out, err := cfg.Render()
	if hasBin {
		assert.NoError(t, err)
		assert.IsType(t, PythonResult{}, out)
		res := out.(PythonResult)
		assert.Empty(t, res.Version)
		assert.Empty(t, res.VersionPy2)
		assert.Empty(t, res.VersionPy3)
		assert.NotEmpty(t, res.Deps)

	} else {
		assert.NoError(t, err)
		assert.Contains(t, cliWarnings.Warnings, "python.deps.packages")
		assert.Len(t, cliWarnings.Warnings["python.deps.packages"], 1)
	}
}

func TestPythonConfig_Render_Err(t *testing.T) {
	cfg := PythonConfig{
		Core: []string{"not-an-option"},
	}

	_, err := cfg.Render()
	assert.Error(t, err)
}

func TestPythonConfig_Render_Empty(t *testing.T) {
	cfg := PythonConfig{}

	out, err := cfg.Render()
	assert.NoError(t, err)
	assert.IsType(t, PythonResult{}, out)

	res := out.(PythonResult)
	assert.True(t, res.IsEmpty())
}

func TestNewPythonResult(t *testing.T) {
	r := NewPythonResult()
	assert.Empty(t, r.Deps)
	assert.Empty(t, r.Version)
	assert.Empty(t, r.VersionPy2)
	assert.Empty(t, r.VersionPy3)
}

func TestPythonResult_IsEmpty(t *testing.T) {
	r := NewPythonResult()
	assert.True(t, r.IsEmpty())
}

func TestPythonResult_IsEmpty2(t *testing.T) {
	r := NewPythonResult()
	r.Version = "3.6"
	assert.False(t, r.IsEmpty())
}

func TestPythonResult_IsEmpty3(t *testing.T) {
	r := NewPythonResult()
	r.VersionPy2 = "2.7"
	assert.False(t, r.IsEmpty())
}

func TestPythonResult_IsEmpty4(t *testing.T) {
	r := NewPythonResult()
	r.VersionPy3 = "3.6"
	assert.False(t, r.IsEmpty())
}

func TestPythonResult_IsEmpty5(t *testing.T) {
	r := NewPythonResult()
	r.Deps["foo"] = "bar"
	assert.False(t, r.IsEmpty())
}

func TestPythonResult_Markdown(t *testing.T) {
	r := NewPythonResult()
	r.Version = "3.6.9"
	r.VersionPy2 = "2.7.11"
	r.VersionPy3 = "3.6.9"
	r.Deps = map[string]string{
		"foo": "1.2.3",
	}

	data, err := r.Markdown()
	assert.NoError(t, err)

	expected := "**Python**\n- _version_: 3.6.9\n- _py2_: 2.7.11\n- _py3_: 3.6.9\n- _dependencies_:\n  ```\n  foo==1.2.3\n  ```\n"
	assert.Equal(t, expected, string(data))
}

func TestPythonResult_Markdown_Empty(t *testing.T) {
	r := NewPythonResult()

	data, err := r.Markdown()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestPythonResult_Plaintext(t *testing.T) {
	r := NewPythonResult()
	r.Version = "3.6.9"
	r.VersionPy2 = "2.7.11"
	r.VersionPy3 = "3.6.9"
	r.Deps = map[string]string{
		"foo": "1.2.3",
	}

	data, err := r.Plaintext()
	assert.NoError(t, err)

	expected := "Python\n------\nversion:  3.6.9\npy2:      2.7.11\npy3:      3.6.9\ndependencies:\n- foo==1.2.3\n"
	assert.Equal(t, expected, string(data))
}

func TestPythonResult_Plaintext_Empty(t *testing.T) {
	r := NewPythonResult()

	data, err := r.Plaintext()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestPythonResult_JSON(t *testing.T) {
	r := NewPythonResult()
	r.Version = "3.6.9"
	r.VersionPy2 = "2.7.11"
	r.VersionPy3 = "3.6.9"
	r.Deps = map[string]string{
		"foo": "1.2.3",
	}

	data, err := r.JSON()
	assert.NoError(t, err)

	expected := `{"version":"3.6.9","py2":"2.7.11","py3":"3.6.9","dependencies":{"foo":"1.2.3"}}`
	assert.Equal(t, expected, string(data))
}

func TestPythonResult_JSON_Empty(t *testing.T) {
	r := NewPythonResult()

	data, err := r.JSON()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

func TestPythonResult_YAML(t *testing.T) {
	r := NewPythonResult()
	r.Version = "3.6.9"
	r.VersionPy2 = "2.7.11"
	r.VersionPy3 = "3.6.9"
	r.Deps = map[string]string{
		"foo": "1.2.3",
	}

	data, err := r.YAML()
	assert.NoError(t, err)

	expected := heredoc.Doc(`
		version: 3.6.9
		py2: 2.7.11
		py3: 3.6.9
		dependencies:
		  foo: 1.2.3
	`)
	assert.Equal(t, expected, string(data))
}

func TestPythonResult_YAML_Empty(t *testing.T) {
	r := NewPythonResult()

	data, err := r.YAML()
	assert.NoError(t, err)
	assert.Empty(t, data)
}

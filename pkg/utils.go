package pkg

import (
	"bytes"
	"os/exec"
	"strings"
)

// binExists is a helper function which checks if a given binary
// exists somewhere on the PATH.
func binExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// normalize is a helper function which strips a given []byte
// of any return characters (\r\n).
func normalize(in []byte) string {
	s := strings.ReplaceAll(string(in), "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// toSlice is a helper function which normalizes a given []byte
// and splits it into string slice using a space as the delimiter.
func toSlice(in []byte) []string {
	s := normalize(in)
	return strings.Split(s, " ")
}

// runCommand is a helper to run a command and collect the output from
// stdout and stderr.
func runCommand(name string, args ...string) (bytes.Buffer, bytes.Buffer, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()

	return stdout, stderr, err
}

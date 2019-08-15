package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinExists_True(t *testing.T) {
	exists := binExists("go")
	assert.True(t, exists)
}

func TestBinExists_False(t *testing.T) {
	exists := binExists("jk3rlkdal3r93")
	assert.False(t, exists)
}

func TestNormalize(t *testing.T) {
	var tests = []struct {
		in  []byte
		out string
	}{
		{[]byte(""), ""},
		{[]byte(" "), " "},
		{[]byte("abc"), "abc"},
		{[]byte("abc def"), "abc def"},
		{[]byte("abc\ndef"), "abcdef"},
		{[]byte("abc\rdef"), "abcdef"},
		{[]byte("abc\r\ndef"), "abcdef"},
		{[]byte("\r\nabc\r\ndef\n\n\n\n"), "abcdef"},
		{[]byte("\r\nabc \r\ndef\n\n\n\n"), "abc def"},
		{[]byte("\r\r\n\r\n\n\n\n\r"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.in), func(t *testing.T) {
			actual := normalize(tt.in)
			assert.Equal(t, tt.out, actual)
		})
	}
}

func TestToSlice(t *testing.T) {
	var tests = []struct {
		in  []byte
		out []string
	}{
		{[]byte(""), []string{""}},
		{[]byte(" "), []string{"", ""}},
		{[]byte("a"), []string{"a"}},
		{[]byte(" a "), []string{"", "a", ""}},
		{[]byte("a b c"), []string{"a", "b", "c"}},
		{[]byte(" a b c "), []string{"", "a", "b", "c", ""}},
		{[]byte("\r\n a\r\n b\r\n c\r\n"), []string{"", "a", "b", "c"}},
		{[]byte("a\nb\nc"), []string{"abc"}},
	}

	for _, tt := range tests {
		t.Run(string(tt.in), func(t *testing.T) {
			actual := toSlice(tt.in)
			assert.Equal(t, tt.out, actual)
		})
	}
}

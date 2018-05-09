package text

import (
	"testing"
)

func TestConvertLineEnding(t *testing.T) {
	cases := []struct {
		str      string
		eol      string
		expected string
	}{
		{"\r", "\n", "\n"},
		{"\n", "\n", "\n"},
		{"\r\n", "\n", "\n"},
		{"\r", "\r\n", "\r\n"},
		{"\n", "\r\n", "\r\n"},
		{"\r\n", "\r\n", "\r\n"},
		{"foo\rbar\r", "\n", "foo\nbar\n"},
		{"foo\nbar\n", "\n", "foo\nbar\n"},
		{"foo\r\nbar\r\n", "\n", "foo\nbar\n"},
		{"foo\nbar\n", "\r\n", "foo\r\nbar\r\n"},
		{"foo\n\nbar\n", "\n", "foo\n\nbar\n"},
		{"foo\r\n\r\nbar\r\n", "\n", "foo\n\nbar\n"},
		{"foo\r\r\nbar\r", "\n", "foo\n\nbar\n"},
		{"foo\n\r\nbar\n", "\n", "foo\n\nbar\n"},
		{"foo\n\r\nbar\n", "\r\n", "foo\r\n\r\nbar\r\n"},
		{"foo\n\r\nbar\n", "", "foo\n\r\nbar\n"},
	}

	for idx, c := range cases {
		actual := ConvertLineEnding(c.str, c.eol)
		if actual != c.expected {
			t.Fatalf("[%d] expect %#v, but %#v", idx, c.expected, actual)

		}
	}
}

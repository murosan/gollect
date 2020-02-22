package gollect

import (
	"testing"
)

func TestTrimQuotes(t *testing.T) {
	cases := []struct{ in, want string }{
		{in: "\"abc\"", want: "abc"},
		{in: "\"\"", want: ""},
		{in: "abc", want: "b"},
	}

	for i, c := range cases {
		if v := trimQuotes(c.in); v != c.want {
			t.Errorf("at: %d, want: %s, actual: %s", i, c.want, v)
		}
	}
}

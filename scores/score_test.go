package score

import "testing"

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  int
	}{{"42", 42}, {" 42\n", 42}, {"-1", 0}, {"bad", 0}, {"", 0}, {"999999999999999999999999999999", 0}} {
		if got := parse(tc.input); got != tc.want {
			t.Errorf("parse(%q)=%d, want %d", tc.input, got, tc.want)
		}
	}
}

package fanuc

import "testing"

func TestStrings(t *testing.T) {
	tests := []struct {
		Type
		Id      int
		Comment string
		Exp     string
	}{
		{Ain, 1, "foo", "AI[1:foo]"},
		{Ain, 2, "", "AI[2]"},
		{Din, 3, "bar", "DI[3:bar]"},
		{Gout, 4, "", "GO[4]"},
	}

	for _, test := range tests {
		io := IO{Type: test.Type, Id: test.Id, Comment: test.Comment}
		if io.String() != test.Exp {
			t.Errorf("Bad string. Got %q, want %q", io, test.Exp)
		}
	}
}

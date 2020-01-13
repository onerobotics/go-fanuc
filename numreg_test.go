package fanuc

import (
	"testing"
)

func TestNumericRegister(t *testing.T) {
	tests := []struct {
		id      int
		comment string
		exp     string
	}{
		{1, "numreg one", "R[1:numreg one]"},
		{2, "", "R[2]"},
	}

	for _, test := range tests {
		n := NumericRegister{Id: test.id, Comment: test.comment}
		if n.Id != test.id {
			t.Errorf("Bad id. Got %d, want %d", n.Id, test.id)
		}
		if n.Comment != test.comment {
			t.Errorf("Bad comment. Got %q, want %q", n.Comment, test.id)
		}
		if n.String() != test.exp {
			t.Errorf("Bad string. Got %q, want %q", n, test.exp)
		}
	}
}

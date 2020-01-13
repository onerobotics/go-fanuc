package fanuc

import (
	"testing"
)

func TestFrameString(t *testing.T) {
	tests := []struct {
		Type
		Id      int
		Comment string
		Exp     string
	}{
		{ToolFrame, 1, "foo", "UTOOL[1:foo]"},
		{ToolFrame, 1, "", "UTOOL[1]"},
		{JogFrame, 2, "", "JOG[2]"},
		{UserFrame, 3, "bar", "UFRAME[3:bar]"},
	}

	for _, test := range tests {
		f := &Frame{Type: test.Type, Id: test.Id, Comment: test.Comment}
		if f.String() != test.Exp {
			t.Errorf("Bad string. Got %q, want %q", f.String(), test.Exp)
		}
	}
}

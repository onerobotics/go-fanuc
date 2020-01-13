package fanuc

import (
	"net/url"
	"testing"
)

func TestNewHTTPClientBadTimeout(t *testing.T) {
	_, err := newHTTPGetter("host", 0)
	exp := "Please specify a timeout > 0"
	if err == nil {
		t.Fatal("want an error")
	}
	if err.Error() != exp {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), exp)
	}
}

func TestNewHTTPClientBadHost(t *testing.T) {
	_, err := newHTTPGetter("\\", 100)
	if _, ok := err.(*url.Error); !ok {
		t.Fatal("wanted an url.Error")
	}
}

/*
func TestNumericRegisters(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer server.Close()

	host := server.URL[7:] // remove http://
	c, err := NewHTTPClient(host, 500)
	if err != nil {
		t.Fatal(err)
	}

	numregs, err := c.NumericRegisters()
	if err != nil {
		t.Fatal(err)
	}

	if len(numregs) != 200 {
		t.Fatalf("Got %d numregs. Want 200", len(numregs))
	}

	tests := []struct {
		index   int
		id      int
		comment string
		value   string
	}{
		{0, 1, "TaskId", "0"},
		{1, 2, "Status", "0"},
		{2, 3, "GripMem", "0"},
		{15, 16, "OutfeedAPLD", "70.000000"},
		{199, 200, "SIM/DryRun", "0"},
	}

	for _, test := range tests {
		n := numregs[test.index]
		if n.Id != test.id {
			t.Errorf("bad id. Got %d, want %d", n.Id, test.id)
		}
		if n.Comment != test.comment {
			t.Errorf("bad comment. Got %q, want %q", n.Comment, test.comment)
		}

		if n.Value != test.value {
			t.Errorf("bad value. Got %q, want %q", n.Value, test.value)
		}
	}
}
*/

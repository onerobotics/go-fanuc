package fanuc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileClient(t *testing.T) {
	c, err := NewFileClient("testdata")
	if err != nil {
		t.Fatal(err)
	}

	allTests(c, t)
}

func TestHTTPClient(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer server.Close()

	c, err := NewHTTPClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	allTests(c, t)
}

func allTests(c Client, t *testing.T) {
	testNumregs(c, t)
	testPosregs(c, t)
	testErrors(c, t)
	testIO(c, t)
}

func testNumregs(c Client, t *testing.T) {
	// numregs
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
		{0, 1, "this is an extre", "0"},
		{1, 2, "two", "0"},
		{2, 3, "three", "0"},
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

func testPosregs(c Client, t *testing.T) {
	posregs, err := c.PositionRegisters()
	if err != nil {
		t.Fatal(err)
	}

	if len(posregs) != 100 {
		t.Fatalf("Got %d posregs. Want 100", len(posregs))
	}

	tests := []struct {
		id      int
		comment string
	}{
		{1, "Maintenance"},
		{2, "pr2"},
		{6, "ZERO"},
		{20, "Outfeed Approach"},
		{97, "M2 Place RTTO"},
	}

	for _, test := range tests {
		pr := posregs[test.id-1]

		if pr.Id != test.id {
			t.Errorf("bad id. Got %d, want %d", pr.Id, test.id)
		}
		if pr.Comment != test.comment {
			t.Errorf("bad comment. Got %q, want %q", pr.Comment, test.comment)
		}
	}

}

func testErrors(c Client, t *testing.T) {
	errors, err := c.Errors()
	if err != nil {
		t.Fatal(err)
	}

	if len(errors) != 99 {
		t.Fatalf("Got %d errors, want 99", len(errors))
	}

	var tests = []struct {
		id       int
		seq      int
		alarm    string
		msg      string
		severity string
	}{
		{2, 1779, "SYST-179", "SHIFT-RESET Released", "WARN"},
		{3, 1780, "SYST-178", "SHIFT-RESET Pressed", "WARN"},
		{14, 1791, "HOST-209", "SM: Connection Aborted", "WARN"},
		{91, 1869, "SYST-043", "TP disabled in T1/T2 mode", "STOP.L"},
	}

	for _, test := range tests {
		err := errors[test.id]

		if err.Sequence != test.seq {
			t.Errorf("bad seq. Got %d, want %d", err.Sequence, test.seq)
		}
		if err.Alarm.String() != test.alarm {
			t.Errorf("Bad alarm. Got %q, want %q", err.Alarm, test.alarm)
		}
		if err.Alarm.Msg != test.msg {
			t.Errorf("Bad msg. Got %q, want %q", err.Alarm.Msg, test.msg)
		}
		if err.Severity != test.severity {
			t.Errorf("Bad severity. Got %q, want %q", err.Severity, test.severity)
		}
	}

}

func testIO(c Client, t *testing.T) {
	tests := []struct {
		Index int
		Type
		Id      int
		Value   string
		Comment string
	}{
		{0, Ain, 1, "0", ""},
		{1, Aout, 1, "0", "test"},
		{2, Gin, 1, "0", "RecipeReadData"},
		{3, Gout, 1, "0", "RecipeEchoData"},
		{4, Uin, 1, "OFF", "*IMSTP"},
		{12, Uin, 9, "OFF", "RSR1/PNS1/STYLE1"},
	}

	io, err := c.IO(Ain, Aout, Gin, Gout, Uin, Uout)
	if err != nil {
		t.Fatal(err)
	}

	if len(io) < len(tests) {
		t.Fatalf("Only got %d signals. Need at least %d", len(io), len(tests))
	}

	for _, test := range tests {
		bit := io[test.Index]

		if bit.Type != test.Type {
			t.Errorf("Bad type. Got %s, want %s", bit.Type, test.Type)
		}
		if bit.Id != test.Id {
			t.Errorf("Bad id. Got %d, want %d", bit.Id, test.Id)
		}
		if bit.Comment != test.Comment {
			t.Errorf("Bad comment. Got %q, want %q", bit.Comment, test.Comment)
		}
	}
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(".")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c.(*FileClient); !ok {
		t.Errorf("Bad type. Got %T, want *FileClient", c)
	}

	c, err = NewClient("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c.(*HTTPClient); !ok {
		t.Errorf("Bad type. Got %T, want *HTTPClient", c)
	}
}

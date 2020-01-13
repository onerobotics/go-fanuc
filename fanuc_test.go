package fanuc

import (
	"testing"
)

func TestFanucStrings(t *testing.T) {
	tests := []struct {
		Type
		Exp     string
		Verbose string
	}{
		{Numreg, "R", "Numeric Register"},
		{Posreg, "PR", "Position Register"},
		{Sreg, "SR", "String Register"},
		{Ualm, "UALM", "User Alarm"},
		{Ain, "AI", "Analog Input"},
		{Aout, "AO", "Analog Output"},
		{Din, "DI", "Digital Input"},
		{Dout, "DO", "Digital Output"},
		{Flag, "F", "Flag"},
		{Gin, "GI", "Group Input"},
		{Gout, "GO", "Group Output"},
		{Rin, "RI", "Robot Input"},
		{Rout, "RO", "Robot Output"},
		{Sin, "SI", "SOP Input"},
		{Sout, "SO", "SOP Output"},
		{Uin, "UI", "UOP Input"},
		{Uout, "UO", "UOP Output"},
	}

	for _, test := range tests {
		if test.Type.String() != test.Exp {
			t.Errorf("Bad string. Got %q, want %q", test.Type.String(), test.Exp)
		}

		if test.Type.VerboseName() != test.Verbose {
			t.Errorf("Bad verbose name. Got %q, want %q", test.Type.VerboseName(), test.Verbose)
		}
	}
}

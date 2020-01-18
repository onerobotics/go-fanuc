package fanuc

import (
	"fmt"
	"strconv"
)

type Type uint

const (
	Invalid Type = iota
	Numreg
	Posreg
	Sreg
	Ualm

	Ain
	Aout
	Din
	Dout
	Flag
	Gin
	Gout
	Rin
	Rout
	Sin
	Sout
	Uin
	Uout

	UserFrame
	ToolFrame
	JogFrame
)

var names = [...]string{
	Numreg: "R",
	Posreg: "PR",
	Sreg:   "SR",
	Ualm:   "UALM",

	Ain:  "AI",
	Aout: "AO",
	Din:  "DI",
	Dout: "DO",
	Flag: "F",
	Gin:  "GI",
	Gout: "GO",
	Rin:  "RI",
	Rout: "RO",
	Sin:  "SI",
	Sout: "SO",
	Uin:  "UI",
	Uout: "UO",

	ToolFrame: "UTOOL",
	UserFrame: "UFRAME",
	JogFrame:  "JOG",
}

var verboseNames = [...]string{
	Numreg: "Numeric Register",
	Posreg: "Position Register",
	Sreg:   "String Register",
	Ualm:   "User Alarm",

	Ain:  "Analog Input",
	Aout: "Analog Output",
	Din:  "Digital Input",
	Dout: "Digital Output",
	Flag: "Flag",
	Gin:  "Group Input",
	Gout: "Group Output",
	Rin:  "Robot Input",
	Rout: "Robot Output",
	Sin:  "SOP Input",
	Sout: "SOP Output",
	Uin:  "UOP Input",
	Uout: "UOP Output",

	ToolFrame: "Tool Frame",
	UserFrame: "User Frame",
	JogFrame:  "Jog Frame",
}

func (t Type) String() string {
	s := ""
	if 0 < t && t < Type(len(names)) {
		s = names[t]
	}
	if s == "" {
		s = "Invalid(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

func (t Type) VerboseName() string {
	s := ""
	if 0 < t && t < Type(len(verboseNames)) {
		s = verboseNames[t]
	}
	if s == "" {
		s = "Invalid(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

type device int

const (
	MD device = iota
	KAREL
)

func (d device) String() string {
	switch d {
	case MD:
		return "MD"
	case KAREL:
		return "KAREL"
	}
	return fmt.Sprintf("device(%d)", d)
}

package fanuc

type commentCode int

const (
	invalid commentCode = iota
	numreg
	_ // set numreg value
	posreg
	ualm
	_ // set ualm_sev value
	rin
	rout
	din
	dout
	gin
	gout
	ain
	aout
	sreg
	_ // set sreg value
	_ // ?
	_ // ?
	_ // ?
	flag
)

func commentCodeFor(t Type) commentCode {
	switch t {
	case Numreg:
		return numreg
	case Posreg:
		return posreg
	case Sreg:
		return sreg
	case Ualm:
		return ualm
	case Ain:
		return ain
	case Aout:
		return aout
	case Din:
		return din
	case Dout:
		return dout
	case Flag:
		return flag
	case Gin:
		return gin
	case Gout:
		return gout
	case Rin:
		return rin
	case Rout:
		return rout
	}

	return invalid
}

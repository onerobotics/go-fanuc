package fanuc

type Rep int

const (
	Joint Rep = iota
	Cartesian
	Uninitialized
)

type Config struct {
	Flip       bool
	Up         bool
	Top        bool
	TurnCounts [3]int
}

type Position struct {
	Id      int
	Comment string

	Group int
	Rep
	Uframe int
	Utool  int
	Config
	X, Y, Z, W, P, R float32
	Joints           []float32
}

package fanuc

import (
	"fmt"
)

// TODO: refactor this ugliness!

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
	TurnCounts []int
}

type PositionRegister struct {
	Id      int
	Comment string

	Group int
	Rep
	Config
	X, Y, Z, W, P, R float32
	Joints           []float32
}

func (p PositionRegister) Value() string {
	switch p.Rep {
	case Uninitialized:
		return fmt.Sprintf(`[%d,%d] =   '%s' Uninitialized`, p.Group, p.Id, p.Comment)
	case Joint:
		return fmt.Sprintf(`[%d,%d] =   '%s'   Group: %d
		J1 =     %f deg   J2 =   %f deg   J3 =   %f deg
		J4 =     %f deg   J5 =   %f deg   J6 =   %f deg`, p.Group, p.Id, p.Comment, p.Group, p.Joints[0], p.Joints[1], p.Joints[2], p.Joints[3], p.Joints[4], p.Joints[5])
	default:
		var cf string
		if p.Config.Flip {
			cf = "F"
		} else {
			cf = "N"
		}
		var cu string
		if p.Config.Up {
			cu = "U"
		} else {
			cu = "D"
		}
		var ct string
		if p.Config.Top {
			ct = "T"
		} else {
			ct = "B"
		}
		return fmt.Sprintf(`[%d,%d] =   '%s'
		Group: %d   Config: %s %s %s, %d, %d, %d
		X:  %f   Y:  %f   Z:  %f
		W:  %f   P:  %f   R:  %f`, p.Group, p.Id, p.Comment, p.Group, cf, cu, ct, p.Config.TurnCounts[0], p.Config.TurnCounts[1], p.Config.TurnCounts[2], p.X, p.Y, p.Z, p.W, p.P, p.R)
	}
}

func (p PositionRegister) String() string {
	comment := p.Comment
	if comment != "" {
		comment = ":" + comment
	}

	return fmt.Sprintf("PR[%d%s]", p.Id, comment)
}

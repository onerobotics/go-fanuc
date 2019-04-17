package fanuc

import (
	"fmt"
	"regexp"
	"strconv"
)

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
	Group int
	Id    int
	Rep
	Config
	X, Y, Z, W, P, R float32
	Joints           []float32
	Comment          string
}

var (
	posregRegex *regexp.Regexp
)

func init() {
	posregRegex = regexp.MustCompile(`(?m)\[(\d),(\d+)\] =   \'([^']*)' (Uninitialized|\r?\n  Group: (\d)   Config: (F|N) (U|D) (T|B), (\d), (\d), (\d)\r?\n  X:\s*(-?\d*.\d+|[*]+)   Y:\s+(-?\d*.\d+|[*]+)   Z:\s+(-?\d*.\d+|[*]+)\r?\n  W:\s*(-?\d*.\d+|[*]+)   P:\s*(-?\d*.\d+|[*]+)   R:\s*(-?\d*.\d+|[*]+)|  Group: (\d)\r?\n  (J1) =\s*(-?\d*.\d+|[*]+) deg   J2 =\s*(-?\d*.\d+|[*]+) deg   J3 =\s*(-?\d*.\d+|[*]+) deg \r?\n  J4 =\s*(-?\d*.\d+|[*]+) deg   J5 =\s*(-?\d*.\d+|[*]+) deg   J6 =\s*(-?\d*.\d+|[*]+) deg)`)
}

func (p *PositionRegister) String() string {
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

func (c *Client) GetPositionRegisters() ([]PositionRegister, error) {
	body, err := c.get("/md/posreg.va")
	if err != nil {
		return nil, err
	}

	var posregs []PositionRegister

	matches := posregRegex.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		group, err := strconv.Atoi(m[1])
		if err != nil {
			return posregs, err
		}

		id, err := strconv.Atoi(m[2])
		if err != nil {
			return posregs, err
		}

		switch {
		case m[4] == "Uninitialized":
			posregs = append(posregs, PositionRegister{
				Group:   group,
				Id:      id,
				Comment: m[3],
				Rep:     Uninitialized,
			})
		case m[16] == "J1":
			j1, _ := strconv.ParseFloat(m[17], 32)
			j2, _ := strconv.ParseFloat(m[18], 32)
			j3, _ := strconv.ParseFloat(m[19], 32)
			j4, _ := strconv.ParseFloat(m[20], 32)
			j5, _ := strconv.ParseFloat(m[21], 32)
			j6, _ := strconv.ParseFloat(m[22], 32)

			posregs = append(posregs, PositionRegister{
				Group:   group,
				Id:      id,
				Comment: m[3],
				Rep:     Joint,
				Joints: []float32{float32(j1),
					float32(j2),
					float32(j3),
					float32(j4),
					float32(j5),
					float32(j6)},
			})
		default:
			cfgFlip := m[6] == "N"
			cfgUp := m[7] == "U"
			cfgTop := m[8] == "T"
			tc1, _ := strconv.Atoi(m[9])
			tc2, _ := strconv.Atoi(m[10])
			tc3, _ := strconv.Atoi(m[11])
			x, _ := strconv.ParseFloat(m[12], 32)
			y, _ := strconv.ParseFloat(m[13], 32)
			z, _ := strconv.ParseFloat(m[14], 32)
			w, _ := strconv.ParseFloat(m[15], 32)
			p, _ := strconv.ParseFloat(m[16], 32)
			r, _ := strconv.ParseFloat(m[17], 32)
			posregs = append(posregs, PositionRegister{
				Group:   group,
				Id:      id,
				Comment: m[3],
				Rep:     Cartesian,
				Config: Config{
					Flip:       cfgFlip,
					Up:         cfgUp,
					Top:        cfgTop,
					TurnCounts: []int{tc1, tc2, tc3},
				},
				X: float32(x),
				Y: float32(y),
				Z: float32(z),
				W: float32(w),
				P: float32(p),
				R: float32(r),
			})
		}
	}

	return posregs, nil
}

package fanuc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TODO: destroy all regexps!

var (
	errorsRegexp      *regexp.Regexp
	framesRegexp      *regexp.Regexp
	ioRegexp          *regexp.Regexp
	numregsRegexp     *regexp.Regexp
	posregsRegexp     *regexp.Regexp
	tpfilenamesRegexp *regexp.Regexp
)

func init() {
	errorsRegexp = regexp.MustCompile(`(\d+)" (\d{2}-\w+-\d{2} \d{2}:\d{2}:\d+) " (R E S E T\s{41}|((\w{4})-(\d{3})) (.{41}))" " (.{30})(\d{8})"(.{4})`)
	framesRegexp = regexp.MustCompile(`(-?\d+\.\d+) *(-?\d+\.\d+) *(-?\d+\.\d+) *(-?\d+\.\d+) *(-?\d+\.\d+) *(-?\d+\.\d+) ([a-zA-Z0-9_ ]*)`)
	ioRegexp = regexp.MustCompile(`(AIN|AOUT|DIN|DOUT|GIN|GOUT|SI|SO|FLG|RI|RO|UI|UO)\[\s*(\d+)\]\s+(ON|OFF|\d+)  ([^\n]{0,24})`)
	numregsRegexp = regexp.MustCompile(`\s+\[(\d+)\] = (-?\d*(\.\d+)?)  '([^']*)'`)
	posregsRegexp = regexp.MustCompile(`(?m)\[(\d),(\d+)\] =   \'([^']*)' (Uninitialized|\r?\n  Group: (\d)   Config: (F|N) (U|D) (T|B), (\d), (\d), (\d)\r?\n  X:\s*(-?\d*.\d+|[*]+)   Y:\s+(-?\d*.\d+|[*]+)   Z:\s+(-?\d*.\d+|[*]+)\r?\n  W:\s*(-?\d*.\d+|[*]+)   P:\s*(-?\d*.\d+|[*]+)   R:\s*(-?\d*.\d+|[*]+)|  Group: (\d)\r?\n  (J1) =\s*(-?\d*.\d+|[*]+) deg   J2 =\s*(-?\d*.\d+|[*]+) deg   J3 =\s*(-?\d*.\d+|[*]+) deg \r?\n  J4 =\s*(-?\d*.\d+|[*]+) deg   J5 =\s*(-?\d*.\d+|[*]+) deg   J6 =\s*(-?\d*.\d+|[*]+) deg)`)
	tpfilenamesRegexp = regexp.MustCompile(`>([A-Z][A-Z0-9_]*)\.TP`)
}

func parseErrors(src string) (errors []Error, err error) {
	matches := errorsRegexp.FindAllStringSubmatch(src, -1)
	for _, m := range matches {
		sequence, err := strconv.Atoi(m[1])
		if err != nil {
			return errors, err
		}

		// 17-APR-19 10:23:08
		timestamp := m[2]
		newTimestamp := timestamp[:7] + "20" + timestamp[7:] // make full year
		t, err := time.Parse("2-Jan-2006 15:04:05", newTimestamp)
		if err != nil {
			return errors, err
		}

		if m[3] == "R E S E T" {
			errors = append(errors, Error{
				Sequence: sequence,
				Time:     t,
			})
			continue
		}

		errors = append(errors, Error{
			Sequence: sequence,
			Time:     t,
			Alarm: Alarm{
				Facility: m[5],
				Code:     m[6],
				Msg:      strings.TrimSpace(m[7]),
			},
			Severity: strings.TrimSpace(m[8]),
		})
	}

	return
}

func parseFrames(src string) (frames []Frame, err error) {
	tf := strings.Index(src, "Tool Frame")
	jf := strings.Index(src, "Jog Frame")
	uf := strings.Index(src, "User Frame")
	if tf < 0 || jf < 0 || uf < 0 {
		return nil, errors.New("Invalid frame.dg")
	}

	tfMatches := framesRegexp.FindAllStringSubmatch(src[tf+10:jf], -1)
	jfMatches := framesRegexp.FindAllStringSubmatch(src[jf+9:uf], -1)
	ufMatches := framesRegexp.FindAllStringSubmatch(src[uf+10:], -1)

	proc := func(matches [][]string, t Type, frames *[]Frame) {
		for id, m := range matches {
			*frames = append(*frames, Frame{
				Id:      id + 1,
				Type:    t,
				X:       m[1],
				Y:       m[2],
				Z:       m[3],
				W:       m[4],
				P:       m[5],
				R:       m[6],
				Comment: m[7],
			})
		}
	}
	proc(tfMatches, ToolFrame, &frames)
	proc(jfMatches, JogFrame, &frames)
	proc(ufMatches, UserFrame, &frames)

	return
}

func parseIO(src string) (io []IO, err error) {
	matches := ioRegexp.FindAllStringSubmatch(src, -1)
	for _, m := range matches {
		id, err := strconv.Atoi(m[2])
		if err != nil {
			return io, err
		}

		var t Type
		switch m[1] {
		case "AIN":
			t = Ain
		case "AOUT":
			t = Aout
		case "DIN":
			t = Din
		case "DOUT":
			t = Dout
		case "FLG":
			t = Flag
		case "GIN":
			t = Gin
		case "GOUT":
			t = Gout
		case "RI":
			t = Rin
		case "RO":
			t = Rout
		case "SI":
			t = Rin
		case "SO":
			t = Rout
		case "UI":
			t = Uin
		case "UO":
			t = Uout
		}

		io = append(io, IO{Type: t, Id: id, Comment: m[4], Value: m[3]})
	}

	return
}

func parseNumericRegisters(src string) (numregs []NumericRegister, err error) {
	matches := numregsRegexp.FindAllStringSubmatch(src, -1)
	for _, m := range matches {
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return numregs, err
		}

		numregs = append(numregs, NumericRegister{Id: id, Value: m[2], Comment: m[4]})
	}

	return
}

func parsePositionRegisters(src string) (posregs []PositionRegister, err error) {
	matches := posregsRegexp.FindAllStringSubmatch(src, -1)
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
				Id:      id,
				Comment: m[3],
				Group:   group,
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
				Id:      id,
				Comment: m[3],

				Group: group,
				Rep:   Joint,
				Joints: []float32{float32(j1),
					float32(j2),
					float32(j3),
					float32(j4),
					float32(j5),
					float32(j6)},
			})
		default:
			cfgFlip := m[6] == "F"
			cfgUp := m[7] == "U"
			cfgTop := m[8] == "T"
			tc1, _ := strconv.Atoi(m[9])
			tc2, _ := strconv.Atoi(m[10])
			tc3, _ := strconv.Atoi(m[11])
			/* TODO: components will have float32 zero value
			 * even if FANUC data is UNINIT (e.g. ****)
			 * should the PositionRegister struct just store
			 * the string values and let the user parse things
			 * out?
			 */
			x, _ := strconv.ParseFloat(m[12], 32)
			y, _ := strconv.ParseFloat(m[13], 32)
			z, _ := strconv.ParseFloat(m[14], 32)
			w, _ := strconv.ParseFloat(m[15], 32)
			p, _ := strconv.ParseFloat(m[16], 32)
			r, _ := strconv.ParseFloat(m[17], 32)
			posregs = append(posregs, PositionRegister{
				Id:      id,
				Comment: m[3],

				Group: group,
				Rep:   Cartesian,
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

	return
}

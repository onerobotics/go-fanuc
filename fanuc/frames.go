package fanuc

import (
	"fmt"
	"regexp"
)

type FrameKind int

const (
	USER FrameKind = iota
	TOOL
	JOG
)

type Frame struct {
	Id int
	FrameKind
	X, Y, Z, W, P, R string
	Comment          string
}

func (f *Frame) String() string {
	kind := "UFRAME"
	if f.FrameKind == TOOL {
		kind = "UTOOL"
	}
	if f.FrameKind == JOG {
		kind = "JOG"
	}
	comment := f.Comment
	if comment != "" {
		comment = ":" + comment
	}

	return fmt.Sprintf("%s[%d%s] %s %s %s %s %s %s", kind, f.Id, comment, f.X, f.Y, f.Z, f.W, f.P, f.R)
}

var (
	frameRegexp *regexp.Regexp
)

func init() {
	frameRegexp = regexp.MustCompile(`(-?\d+\.\d+)\s+(-?\d+\.\d+)\s+(-?\d+\.\d+)\s+(-?\d+\.\d+)\s+(-?\d+\.\d+)\s+(-?\d+\.\d+) ([^\n]*)`)
}

func (c *Client) GetFrames() ([]Frame, error) {
	body, err := c.get("/md/frame.dg")
	if err != nil {
		return nil, err
	}

	var frames []Frame

	matches := frameRegexp.FindAllStringSubmatch(string(body), -1)
	index := 1
	kind := TOOL
	for id, m := range matches {
		if id == 10 {
			kind = JOG
			index = 1
		}
		if id == 16 {
			kind = USER
			index = 1
		}

		frames = append(frames, Frame{
			Id:        index,
			FrameKind: kind,
			X:         m[1],
			Y:         m[2],
			Z:         m[3],
			W:         m[4],
			P:         m[5],
			R:         m[6],
			Comment:   m[7],
		})

		index++
	}

	return frames, nil
}

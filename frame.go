package fanuc

import (
	"fmt"
)

type Frame struct {
	Type
	Id               int
	X, Y, Z, W, P, R string
	Comment          string
}

func (f *Frame) String() string {
	comment := f.Comment
	if comment != "" {
		comment = ":" + comment
	}

	return fmt.Sprintf("%s[%d%s]", f.Type, f.Id, comment)
}

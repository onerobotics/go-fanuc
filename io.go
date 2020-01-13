package fanuc

import (
	"fmt"
)

type IO struct {
	Type
	Id      int
	Comment string
	Value   string
}

func (i *IO) String() string {
	comment := i.Comment
	if comment != "" {
		comment = ":" + comment
	}
	return fmt.Sprintf("%s[%d%s]", i.Type, i.Id, comment)
}

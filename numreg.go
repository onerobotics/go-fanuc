package fanuc

import (
	"fmt"
)

type NumericRegister struct {
	Id      int
	Comment string
	Value   string // TODO: intval, floatval
}

func (n NumericRegister) String() string {
	comment := n.Comment
	if comment != "" {
		comment = ":" + comment
	}

	return fmt.Sprintf("R[%d%s]", n.Id, comment)
}

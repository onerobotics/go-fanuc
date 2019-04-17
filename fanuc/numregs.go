package fanuc

import (
	"fmt"
	"regexp"
	"strconv"
)

type NumericRegister struct {
	Id      int
	Value   string
	Comment string
}

func (n *NumericRegister) String() string {
	comment := n.Comment
	if comment != "" {
		comment = ":" + comment
	}

	return fmt.Sprintf("R[%d%s]  %s", n.Id, comment, n.Value)
}

var (
	numregRegex *regexp.Regexp
)

func init() {
	numregRegex = regexp.MustCompile(`\s+\[(\d+)\] = (-?\d*(\.\d+)?)  '([^']*)'`)
}

func (c *Client) GetNumericRegisters() ([]NumericRegister, error) {
	body, err := c.get("/md/numreg.va")
	if err != nil {
		return nil, err
	}

	var numregs []NumericRegister

	matches := numregRegex.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return numregs, err
		}

		numregs = append(numregs, NumericRegister{id, m[2], m[4]})
	}

	return numregs, nil
}

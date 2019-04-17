package fanuc

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	tpFilenameRegex *regexp.Regexp
)

func init() {
	tpFilenameRegex = regexp.MustCompile(`>([A-Z][A-Z0-9_]*)\.TP`)
}

func (c *Client) GetTPProgramNames() ([]string, error) {
	body, err := c.get("/md/index_tp.htm?_TEMPLATE=")
	if err != nil {
		return nil, err
	}

	var names []string

	matches := tpFilenameRegex.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		names = append(names, m[1])
	}

	return names, nil
}

func (c *Client) GetTPProgramSource(name string) (string, error) {
	body, err := c.get(fmt.Sprintf("/md/%s.LS?_TEMPLATE=", name))
	if err != nil {
		return "", err
	}

	// remove the HTML prefix with [88:]
	src := strings.Split(string(body)[88:], "</XMP>")[0]

	return src, nil
}

package fanuc

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	ioRegex *regexp.Regexp
)

func init() {
	ioRegex = regexp.MustCompile(`(DIN|DOUT|SI|SO|FLG|RI|RO)\[\s*(\d+)\]\s+(ON|OFF)  ([^\n]{0,24})`)

	keywords = make(map[string]IOKind)
	for i := 0; i < len(kinds); i++ {
		keywords[kinds[i]] = IOKind(i)
	}
}

type IOKind int

const (
	INVALID IOKind = iota
	DIN
	DOUT
	SIN
	SOUT
	FLAG
	RIN
	ROUT
	GIN
	GOUT
)

var kinds = [...]string{
	INVALID: "INVALID",
	DIN:     "DIN",
	DOUT:    "DOUT",
	SIN:     "SI",
	SOUT:    "SO",
	FLAG:    "FLG",
	RIN:     "RI",
	ROUT:    "RO",
	GIN:     "GIN",
	GOUT:    "GOUT",
}

func (k IOKind) String() string {
	s := ""
	if 0 <= k && k < IOKind(len(kinds)) {
		s = kinds[k]
	}
	if s == "" {
		s = "kind(" + strconv.Itoa(int(k)) + ")"
	}
	return s
}

var keywords map[string]IOKind

func Lookup(ident string) IOKind {
	if k, is_kind := keywords[ident]; is_kind {
		return k
	}
	return INVALID
}

type Port struct {
	Id int
	IOKind
	State   string
	Comment string
}

func (p *Port) String() string {
	return fmt.Sprintf("%s[%d]  %s  %s", p.IOKind, p.Id, p.State, p.Comment)
}

func (c *Client) GetIOState() ([]Port, error) {
	body, err := c.get("/md/iostate.dg")
	if err != nil {
		return nil, err
	}

	var ports []Port

	matches := ioRegex.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		id, err := strconv.Atoi(m[2])
		if err != nil {
			return ports, err
		}

		kind := Lookup(m[1])

		ports = append(ports, Port{
			Id:      id,
			IOKind:  kind,
			State:   m[3],
			Comment: m[4],
		})
	}

	return ports, nil
}

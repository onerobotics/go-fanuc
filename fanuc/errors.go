package fanuc

import (
	"regexp"
	"strconv"
	"time"
)

type Alarm struct {
	Facility string
	Code     string
	Msg      string
}

type Error struct {
	Sequence int
	Time     time.Time
	Alarm    // nil indicates RESET
	Severity string
}

var (
	errorRegex *regexp.Regexp
)

func init() {
	errorRegex = regexp.MustCompile(`(\d+)" (\d{2}-\w+-\d{2} \d{2}:\d{2}:\d+) " (R E S E T\s{41}|((\w{4})-(\d{3})) (.{41}))" " (.{30})(\d{8})"(.{4})`)
}

func (c *Client) GetErrors() ([]Error, error) {
	body, err := c.get("/md/errall.ls")
	if err != nil {
		return nil, err
	}

	var errors []Error
	matches := errorRegex.FindAllStringSubmatch(string(body), -1)
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
				Msg:      m[7],
			},
			Severity: m[9],
		})
	}

	return errors, nil
}

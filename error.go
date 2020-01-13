package fanuc

import (
	"fmt"
	"time"
)

type Alarm struct {
	Facility string
	Code     string
	Msg      string
}

func (a *Alarm) String() string {
	return fmt.Sprintf("%s-%s", a.Facility, a.Code)
}

type Error struct {
	Sequence int
	Time     time.Time
	Alarm    // nil indicates RESET
	Severity string
}

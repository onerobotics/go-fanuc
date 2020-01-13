package fanuc

import (
	"testing"
)

func TestAlarmString(t *testing.T) {
	a := &Alarm{Facility: "FOO", Code: "123"}
	if a.String() != "FOO-123" {
		t.Errorf("bad string. Got %q, want %q", a.String(), "FOO-123")
	}
}

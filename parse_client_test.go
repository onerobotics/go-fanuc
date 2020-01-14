package fanuc

import (
	"testing"
)

func getFunc(filename string) (string, error) {
	return `1779" 11-JAN-20 21:40:16 " SYST-179 SHIFT-RESET Released                     " " WARN                          00000000"    "`, nil
}

func TestParseClient(t *testing.T) {
	c := &ParseClient{GetFunc: getFunc}
	errs, err := c.Errors()
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 1 {
		t.Fatalf("Got %d errors. Want 1.", len(errs))
	}
}

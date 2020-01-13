package fanuc

import (
	"testing"
)

func TestNewFileGetter(t *testing.T) {
	_, err := newFileGetter("foobar")
	if err == nil {
		t.Fatal("Wanted an error. Got none")
	}
	want := "\"foobar\" does not exist"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	_, err = newFileGetter("file_getter_test.go")
	if err == nil {
		t.Fatal("Wanted an error. Got none")
	}
	want = "\"file_getter_test.go\" is not a directory"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}
}

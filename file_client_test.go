package fanuc

import (
	"testing"
)

func TestNewFileClient(t *testing.T) {
	_, err := NewFileClient("foobar")
	if err == nil {
		t.Fatal("Wanted an error. Got none")
	}
	want := "\"foobar\" does not exist"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}

	_, err = NewFileClient("file_client_test.go")
	if err == nil {
		t.Fatal("Wanted an error. Got none")
	}
	want = "\"file_client_test.go\" is not a directory"
	if err.Error() != want {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), want)
	}
}

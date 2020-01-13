package fanuc

import (
	"net/url"
	"testing"
)

func TestNewHTTPClientBadTimeout(t *testing.T) {
	_, err := newHTTPGetter("host", 0)
	exp := "Please specify a timeout > 0"
	if err == nil {
		t.Fatal("want an error")
	}
	if err.Error() != exp {
		t.Errorf("Bad error msg. Got %q, want %q", err.Error(), exp)
	}
}

func TestNewHTTPClientBadHost(t *testing.T) {
	_, err := newHTTPGetter("\\", 100)
	if _, ok := err.(*url.Error); !ok {
		t.Fatal("wanted an url.Error")
	}
}

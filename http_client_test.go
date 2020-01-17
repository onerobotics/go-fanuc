package fanuc

import (
	"net/url"
	"testing"
)

func TestSetTimeout(t *testing.T) {
	c, err := NewHTTPClient("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	c.SetTimeout(42)
	if c.client.Timeout != 42 {
		t.Errorf("Bad custom timeout. Got %d, want %d", c.client.Timeout, 42)
	}
}

func TestNewHTTPClientBadHost(t *testing.T) {
	_, err := NewHTTPClient("\\")
	if _, ok := err.(*url.Error); !ok {
		t.Fatal("wanted an url.Error")
	}
}

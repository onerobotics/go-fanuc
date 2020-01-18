package fanuc

import (
	"net/url"
	"testing"
)

func TestNewHTTPClient(t *testing.T) {
	goodHosts := []string{"127.0.0.1", "10.0.0.1", "http://foo.bar", "https://doubtfully.secure.robot", "http://127.0.0.1:8080"}
	for _, host := range goodHosts {
		_, err := NewHTTPClient(host)
		if err != nil {
			t.Errorf("NewHTTPClient(%s) failed: %s", host, err)
		}
	}

	badHosts := []string{"foobar", "www.google.com", "http//127.0.0.1"}
	for _, host := range badHosts {
		_, err := NewHTTPClient(host)
		if err == nil {
			t.Errorf("NewHTTPClient(%s) should have failed", host)
		}
	}
}

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
	_, err := NewHTTPClient("foobar")
	if err == nil {
		t.Fatal("Wanted an error")
	}
	if _, ok := err.(*url.Error); !ok {
		t.Fatal("wanted an url.Error")
	}
}

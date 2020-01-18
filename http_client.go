package fanuc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrForbidden    = errForbidden()
	ErrUnauthorized = errUnauthorized()
)

func errForbidden() error {
	return errors.New("Forbidden. Please unlock KAREL via Setup > Host Comm > HTTP.")
}

func errUnauthorized() error {
	return errors.New("Unauthorized. Please unlock KAREL via Setup > Host Comm > HTTP.")
}

type HTTPClient struct {
	client  *http.Client
	baseURL *url.URL

	base *ParseClient
}

const (
	// units are in seconds
	DefaultClientTimeout    = 5
	DefaultDialTimeout      = 1
	DefaultHandshakeTimeout = 1
	DefaultKeepAlive        = 30
)

func NewHTTPClient(s string) (*HTTPClient, error) {
	if ip := net.ParseIP(s); ip != nil {
		s = "http://" + s
	}
	baseURL, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, err
	}

	c := &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: DefaultClientTimeout * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   DefaultDialTimeout * time.Second,
					KeepAlive: DefaultKeepAlive * time.Second,
				}).DialContext,
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: DefaultHandshakeTimeout * time.Second,
			},
		},
	}

	c.base = &ParseClient{GetFunc: c.get}

	return c, nil
}

func (c *HTTPClient) SetTimeout(t time.Duration) {
	c.client.Timeout = t
}

func (c *HTTPClient) get(dev device, path string) (string, error) {
	url, err := c.baseURL.Parse(dev.String() + "/" + path)
	if err != nil {
		return "", err
	}

	params := url.Query()
	params.Add("_TEMPLATE", "")
	url.RawQuery = params.Encode()

	s := strings.ReplaceAll(url.String(), "+", "%20") // FANUC likes %20 over + for spaces

	resp, err := c.client.Get(s)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *HTTPClient) Errors() ([]Error, error) {
	return c.base.Errors()
}

func (c *HTTPClient) Frames() ([]Frame, error) {
	return c.base.Frames()
}

func (c *HTTPClient) IO(types ...Type) ([]IO, error) {
	return c.base.IO(types...)
}

func (c *HTTPClient) NumericRegisters() ([]NumericRegister, error) {
	return c.base.NumericRegisters()
}

func (c *HTTPClient) PositionRegisters() ([]PositionRegister, error) {
	return c.base.PositionRegisters()
}

func (c *HTTPClient) SetComment(t Type, id int, comment string) error {
	fc := commentCodeFor(t)
	if fc == invalid {
		return fmt.Errorf("cannot set comment for %s", t.VerboseName())
	}

	rel := fmt.Sprintf("ComSet?sFc=%d&sIndx=%d&sComment=%s", fc, id, comment)
	_, err := c.get(KAREL, rel)
	return err
}

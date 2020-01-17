package fanuc

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

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

func NewHTTPClient(host string) (*HTTPClient, error) {
	baseURL, err := url.Parse(fmt.Sprintf("http://%s/md/", host))
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

func (c *HTTPClient) get(path string) (result string, err error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return
	}

	resp, err := c.client.Get(u.String() + "?_TEMPLATE=") // remove HTML BS
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(body), nil
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

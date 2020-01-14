package fanuc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL *url.URL

	base *ParseClient
}

func NewHTTPClient(host string, timeout int) (*HTTPClient, error) {
	baseURL, err := url.Parse(fmt.Sprintf("http://%s/md/", host))
	if err != nil {
		return nil, err
	}

	if timeout <= 0 {
		return nil, errors.New("timeout must be > 0")
	}

	c := &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
	}
	c.base = &ParseClient{GetFunc: c.get}

	return c, nil
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

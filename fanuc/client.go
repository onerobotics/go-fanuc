package fanuc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL *url.URL
}

func NewClient(host string) (*Client, error) {
	baseURL, err := url.Parse(fmt.Sprintf("http://%s/", host))
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL: baseURL,
	}, nil
}

func (c *Client) get(urlStr string) ([]byte, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

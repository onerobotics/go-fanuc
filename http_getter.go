package fanuc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type httpGetter struct {
	client  *http.Client
	baseURL *url.URL
}

func newHTTPGetter(host string, timeout int) (*httpGetter, error) {
	if timeout <= 0 {
		return nil, errors.New("Please specify a timeout > 0")
	}

	baseURL, err := url.Parse(fmt.Sprintf("http://%s/md/", host))
	if err != nil {
		return nil, err
	}

	return &httpGetter{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
	}, nil
}

func (c *httpGetter) Get(path string) (result string, err error) {
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

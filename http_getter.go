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

func (c *httpGetter) get(path string) (result string, err error) {
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

/*
func (c *Client) Errors() ([]fanuc.Error, error) {
	body, err := c.get("errall.ls")
	if err != nil {
		return nil, err
	}

	return regexp.Errors(body)
}

func (c *Client) Frames() ([]fanuc.Frame, error) {
	body, err := c.get("frame.dg")
	if err != nil {
		return nil, err
	}

	return regexp.Frames(body)
}

func (c *Client) IO(portType fanuc.Type) ([]fanuc.Port, error) {
	// try to populate portsCache if necessary
	if c.portsCache == nil {
		_, err := c.Ports()
		if err != nil {
			return nil, err
		}
	}

	// only return ports of this type
	var ports []fanuc.Port
	for _, port := range c.portsCache {
		if port.Type == portType {
			ports = append(ports, port)
		}
	}
	return ports, nil
}

func (c *Client) Ports() ([]fanuc.Port, error) {
	body, err := c.get("iostate.dg")
	if err != nil {
		return nil, err
	}

	ports, err := regexp.Ports(body)
	if err != nil {
		return nil, err
	}

	c.portsCache = ports

	return ports, nil
}

func (c *Client) NumericRegisters() ([]fanuc.NumericRegister, error) {
	body, err := c.get("numreg.va")
	if err != nil {
		return nil, err
	}

	return regexp.NumericRegisters(body)
}

func (c *Client) PositionRegisters() ([]fanuc.PositionRegister, error) {
	body, err := c.get("posreg.va")
	if err != nil {
		return nil, err
	}

	return regexp.PositionRegisters(body)
}

/*
func (c *Client) TPPrograms() ([]fanuc.TPProgram, error) {
	body, err := c.get("index_tp.htm")
	if err != nil {
		return nil, err
	}

	var progs []TPProgram
	matches := tpFilenameRegex.FindAllStringSubmatch(string(body), -1)
	for _, m := range matches {
		progs = append(progs, TPProgram{Name: m[1]})
	}

	return progs, nil
}

func (c *Client) TPSource(name string) (src string, err error) {
	body, err := c.get(fmt.Sprintf("%s.LS?_TEMPLATE=", name))
	if err != nil {
		return
	}

	// remove the HTML prefix with [88:]
	// TODO: this may change if FANUC changes their headers
	src = strings.Split(string(body)[88:], "</XMP>")[0]
	return
}
*/

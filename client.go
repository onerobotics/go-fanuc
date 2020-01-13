package fanuc

type Getter interface {
	Get(string) (string, error)
}

type Client struct {
	Getter

	io []IO
}

func NewFileClient(dir string) (*Client, error) {
	g, err := newFileGetter(dir)
	if err != nil {
		return nil, err
	}

	return &Client{Getter: g}, nil
}

func NewHTTPClient(host string, timeout int) (*Client, error) {
	g, err := newHTTPGetter(host, timeout)
	if err != nil {
		return nil, err
	}

	return &Client{Getter: g}, nil
}

func (c *Client) Errors() ([]Error, error) {
	body, err := c.Get("errall.ls")
	if err != nil {
		return nil, err
	}

	return parseErrors(body)
}

func (c *Client) Frames() ([]Frame, error) {
	body, err := c.Get("frame.dg")
	if err != nil {
		return nil, err
	}

	return parseFrames(body)
}

func (c *Client) cacheIO() error {
	body, err := c.Get("iostate.dg")
	if err != nil {
		return err
	}

	io, err := parseIO(body)
	if err != nil {
		return err
	}

	c.io = io

	return nil
}

func contains(t Type, types []Type) bool {
	for _, i := range types {
		if i == t {
			return true
		}
	}

	return false
}

func (c *Client) IO(types ...Type) ([]IO, error) {
	if c.io == nil {
		err := c.cacheIO()
		if err != nil {
			return nil, err
		}
	}

	// Return all IO if no types provided
	if len(types) == 0 {
		return c.io, nil
	}

	var result []IO
	for _, port := range c.io {
		if contains(port.Type, types) {
			result = append(result, port)
		}
	}
	return result, nil
}

func (c *Client) NumericRegisters() ([]NumericRegister, error) {
	body, err := c.Get("numreg.va")
	if err != nil {
		return nil, err
	}

	return parseNumericRegisters(body)
}

func (c *Client) PositionRegisters() ([]PositionRegister, error) {
	body, err := c.Get("posreg.va")
	if err != nil {
		return nil, err
	}

	return parsePositionRegisters(body)
}

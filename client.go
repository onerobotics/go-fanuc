package fanuc

type getter interface {
	get(string) (string, error)
}

type client struct {
	getter

	io []IO
}

func NewFileClient(dir string) (*client, error) {
	g, err := newFileGetter(dir)
	if err != nil {
		return nil, err
	}

	return &client{getter: g}, nil
}

func NewHTTPClient(host string, timeout int) (*client, error) {
	g, err := newHTTPGetter(host, timeout)
	if err != nil {
		return nil, err
	}

	return &client{getter: g}, nil
}

func (c *client) Errors() ([]Error, error) {
	body, err := c.get("errall.ls")
	if err != nil {
		return nil, err
	}

	return parseErrors(body)
}

func (c *client) Frames() ([]Frame, error) {
	body, err := c.get("frame.dg")
	if err != nil {
		return nil, err
	}

	return parseFrames(body)
}

func (c *client) cacheIO() error {
	body, err := c.get("iostate.dg")
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

func (c *client) IO(types ...Type) ([]IO, error) {
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

func (c *client) NumericRegisters() ([]NumericRegister, error) {
	body, err := c.get("numreg.va")
	if err != nil {
		return nil, err
	}

	return parseNumericRegisters(body)
}

func (c *client) PositionRegisters() ([]PositionRegister, error) {
	body, err := c.get("posreg.va")
	if err != nil {
		return nil, err
	}

	return parsePositionRegisters(body)
}

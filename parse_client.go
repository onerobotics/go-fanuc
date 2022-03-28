package fanuc

import (
	"strings"
)

type GetFunc = func(device, string) (string, error)

type ParseClient struct {
	GetFunc

	io []IO
}

func (c *ParseClient) Errors() ([]Error, error) {
	body, err := c.GetFunc(MD, "errall.ls")
	if err != nil {
		return nil, err
	}

	return parseErrors(body)
}

func (c *ParseClient) Frames() ([]Frame, error) {
	body, err := c.GetFunc(MD, "frame.dg")
	if err != nil {
		return nil, err
	}

	return parseFrames(body)
}

func (c *ParseClient) cacheIO() error {
	body, err := c.GetFunc(MD, "iostate.dg")
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

func (c *ParseClient) IO(types ...Type) ([]IO, error) {
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

func (c *ParseClient) NumericRegisters() ([]NumericRegister, error) {
	body, err := c.GetFunc(MD, "numreg.va")
	if err != nil {
		return nil, err
	}

	return parseNumericRegisters(body)
}

func (c *ParseClient) PositionRegisters() ([]PositionRegister, error) {
	body, err := c.GetFunc(MD, "posreg.va")
	if err != nil {
		return nil, err
	}

	return parsePositionRegisters(body)
}

func (c *ParseClient) TPPrograms() ([]string, error) {
	body, err := c.GetFunc(MD, "index_tp.htm")
	if err != nil {
		return nil, err
	}

	return parseTPPrograms(body)
}

func (c *ParseClient) TPPositions(programName string) ([]Position, error) {
	lsName := strings.Replace(strings.ToLower(programName), ".tp", ".ls", -1)
	body, err := c.GetFunc(MD, lsName)
	if err != nil {
		return nil, err
	}

	return parseTPPositions(body)
}

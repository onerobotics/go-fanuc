package fanuc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type FileClient struct {
	dir string

	base *ParseClient
}

func NewFileClient(dir string) (*FileClient, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("%q does not exist", dir)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}

	c := &FileClient{dir: dir}
	c.base = &ParseClient{GetFunc: c.get}

	return c, nil
}

func (c *FileClient) get(dev device, filename string) (string, error) {
	// TODO: support other devices
	if dev != MD {
		return "", fmt.Errorf("%s device not supported", dev)
	}

	b, err := ioutil.ReadFile(path.Join(c.dir, filename))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *FileClient) Errors() ([]Error, error) {
	return c.base.Errors()
}

func (c *FileClient) Frames() ([]Frame, error) {
	return c.base.Frames()
}

func (c *FileClient) IO(types ...Type) ([]IO, error) {
	return c.base.IO(types...)
}

func (c *FileClient) NumericRegisters() ([]NumericRegister, error) {
	return c.base.NumericRegisters()
}

func (c *FileClient) PositionRegisters() ([]PositionRegister, error) {
	return c.base.PositionRegisters()
}

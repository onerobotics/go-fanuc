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

func (c *FileClient) get(filename string) (result string, err error) {
	body, err := ioutil.ReadFile(path.Join(c.dir, filename))
	if err != nil {
		return
	}

	return string(body), nil
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

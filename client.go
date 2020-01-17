package fanuc

import (
	"fmt"
	"net"
	"os"
)

type ErrorClient interface {
	Errors() ([]Error, error)
}

type FrameClient interface {
	Frames() ([]Frame, error)
}

type IOClient interface {
	IO(...Type) ([]IO, error)
}

type NumericRegisterClient interface {
	NumericRegisters() ([]NumericRegister, error)
}

type PositionRegisterClient interface {
	PositionRegisters() ([]PositionRegister, error)
}

type Client interface {
	ErrorClient
	FrameClient
	IOClient
	NumericRegisterClient
	PositionRegisterClient
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}

	return false
}

func isIP(path string) bool {
	ip := net.ParseIP(path)
	if ip == nil {
		return false
	}

	return true
}

func NewClient(target string) (Client, error) {
	switch {
	case isDir(target):
		c, err := NewFileClient(target)
		if err != nil {
			return nil, err
		}

		return c, nil
	case isIP(target):
		c, err := NewHTTPClient(target)
		if err != nil {
			return nil, err
		}

		return c, nil
	default:
		return nil, fmt.Errorf("%q is not a valid IP address or directory", target)
	}
}

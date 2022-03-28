package fanuc

import (
	"fmt"
	"net"
	"net/url"
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
	TPPrograms() ([]string, error)
	TPPositions(programName string) ([]Position, error)
}

func isDir(s string) bool {
	info, err := os.Stat(s)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

func isIP(s string) bool {
	ip := net.ParseIP(s)
	if ip != nil {
		return true
	}
	return false
}

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}

// NewClient returns a *FileClient or *HTTPClient depending on the provided
// argument. Local directories will return a *FileClient. Valid IP addresses
// or valid URL return an *HTTPClient.
// on the provided argument. Local directories will return a
func NewClient(target string) (Client, error) {
	switch {
	case isDir(target):
		c, err := NewFileClient(target)
		if err != nil {
			return nil, err
		}

		return c, nil
	case isIP(target):
		url := fmt.Sprintf("http://%s", target)
		c, err := NewHTTPClient(url)
		if err != nil {
			return nil, err
		}

		return c, nil
	case isURL(target):
		c, err := NewHTTPClient(target)
		if err != nil {
			return nil, err
		}

		return c, nil
	default:
		return nil, fmt.Errorf("%q is not a valid directory or host", target)
	}
}

package fanuc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type fileGetter struct {
	dir string
}

func newFileGetter(dir string) (*fileGetter, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("%q does not exist", dir)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}

	return &fileGetter{dir: dir}, nil
}

func (c *fileGetter) get(filename string) (result string, err error) {
	body, err := ioutil.ReadFile(path.Join(c.dir, filename))
	if err != nil {
		return
	}

	return string(body), nil
}

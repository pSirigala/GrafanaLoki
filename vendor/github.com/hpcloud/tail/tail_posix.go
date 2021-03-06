// +build linux darwin freebsd netbsd openbsd

package tail

import (
	"os"
	"path/filepath"
)

func OpenFile(name string) (file *os.File, err error) {
	filename := name
	// Check if the path requested is a symbolic link
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		filename, err = filepath.EvalSymlinks(name)
		if err != nil {
			return nil, err
		}
	}
	return os.Open(filename)
}

package files

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func GetFilesList(path string, address string) ([]string, error) {
	const op = "internal.lib.GetFilesList"
	dir, err := os.Open(path)
	if err != nil {
		return nil, errors.Errorf("%s - %s", op, err)
	}
	defer dir.Close()
	files, err := dir.Readdirnames(-1)
	for i, _ := range files {
		files[i] = "/" + filepath.Join(path, files[i])
	}
	if err != nil {
		return nil, errors.Errorf("%s - %s", op, err)
	}
	return files, nil
}

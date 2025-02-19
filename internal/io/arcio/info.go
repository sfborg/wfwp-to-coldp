package arcio

import (
	"os"
	"path/filepath"

	"wfwp-to-coldp/internal/ent/wfwp"
	"github.com/gnames/gnsys"
)

func (a *arcio) DirInfo() error {
	// get all fiels with their paths inside of the archive.
	paths, err := a.getFiles()
	if err != nil {
		return err
	}

	// find directory where data files reside
	dataDir := getDataDir(paths)
	var dt wfwp.DataType

	for _, v := range paths {
		dir, file, _ := gnsys.SplitPath(v)

		if dir != dataDir {
			continue
		}

		dt = wfwp.NewDataType(file)
		if dt != wfwp.UnkownDT {
			a.dataPaths[dt] = v
		}
	}

	return nil
}

// getDataDir returns dataDir where data files are residing.
func getDataDir(paths []string) string {
	dirs := make(map[string]int)
	for _, v := range paths {
		dir, file, _ := gnsys.SplitPath(v)
		if wfwp.NewDataType(file) != wfwp.UnkownDT {
			dirs[dir]++
		}
	}

	var res string
	var count int
	for k, v := range dirs {
		if v > count {
			count = v
			res = k
		}
	}
	return res
}

func (a *arcio) getFiles() ([]string, error) {
	var files []string
	root := a.cfg.ExtractDir

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err

	}
	return files, nil
}
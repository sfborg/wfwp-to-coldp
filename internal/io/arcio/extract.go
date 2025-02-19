package arcio

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"wfwp-to-coldp/internal/ent/wfwp"
	"github.com/gnames/gnsys"
)

func (a *arcio) Extract() error {
	var err error
	if strings.HasPrefix(a.path, "http") {
		a.path, err = gnsys.Download(a.path, a.cfg.DownloadDir, true)
		if err != nil {
			return err
		}
	}
	err = a.unzip()
	if err != nil {
		return err
	}
	return nil
}

func (a *arcio) unzip() error {
	exists, _ := gnsys.FileExists(a.path)
	if !exists {
		return &wfwp.ErrorFileMissing{Path: a.path}
	}

	// Open the zip file for reading.
	r, err := zip.OpenReader(a.path)
	if err != nil {
		return &wfwp.ErrExtract{Path: a.path, Err: err}
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the full path for the file/directory and ensure its directory exists.
		fpath := filepath.Join(a.cfg.ExtractDir, f.Name)
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return &wfwp.ErrExtract{Path: fpath, Err: err}
		}

		// If it's a directory, move on to the next entry.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Open the file within the zip.
		rc, err := f.Open()
		if err != nil {
			return &wfwp.ErrExtract{Path: fpath, Err: err}
		}
		defer rc.Close()

		// Create a file in the filesystem.
		outFile, err := os.OpenFile(
			fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode(),
		)
		if err != nil {
			return &wfwp.ErrExtract{Path: fpath, Err: err}
		}
		defer outFile.Close()

		// Copy the contents of the file from the zip to the new file.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return &wfwp.ErrExtract{Path: fpath, Err: err}
		}
	}

	return nil
}
package schemaio

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnames/gnsys"
	"wfwp-to-coldp/internal/ent/wfwplib"
)

type schemaio struct {
	path string
}

func New(path string) wfwplib.Schema {
	res := &schemaio{path: path}
	return res
}

func (s *schemaio) Fetch() ([]byte, error) {
	res, err := s.loadSchema()
	if err == nil {
		return res, nil
	}

	res, err = s.loadSchema()
	return res, err
}

// Clean removes WF/WP data directory.
func (s *schemaio) Clean() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		return &wfwplib.ErrDirRemove{Dir: s.path, Err: err}
	}
	return nil
}

// Path returns temporary path where WF/WP schema is stored.
func (s *schemaio) Path() string {
	return s.path
}

func (s *schemaio) loadSchema() ([]byte, error) {
	var err error
	var exists bool
	schemaPath := filepath.Join(s.path, "schema.sql")

	var schemaSQL = `
DROP TABLE IF EXISTS wfwp;
CREATE TABLE wfwp (
	id INTEGER PRIMARY KEY,
	taxon TEXT,
	number TEXT,
	reference TEXT,
	common_name TEXT,
	distribution TEXT,
	synonyms TEXT,
	status TEXT,
	remarks TEXT,
	conservation_status TEXT,
	photo TEXT,
	orientation TEXT,
	photographer TEXT
);
`

	
	err = os.WriteFile(schemaPath, []byte(schemaSQL), 0644)
	if err != nil {
		err = fmt.Errorf("cannot write %s: %w", schemaPath, err)
		return nil, err
	}

	exists, err = gnsys.FileExists(schemaPath)

	if err != nil {
		err = fmt.Errorf("bad file %s: %w", schemaPath, err)
		return nil, err
	}

	if !exists {
		err = fmt.Errorf("file %s does not exist", schemaPath)
		return nil, err
	}

	res, err := os.ReadFile(schemaPath)
	if err != nil {
		err = fmt.Errorf("cannot read %s: %w", schemaPath, err)
		return nil, err
	}

	return res, nil
}

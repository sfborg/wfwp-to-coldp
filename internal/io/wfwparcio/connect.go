package wfwparcio

import (
	"log/slog"
	"os"
	"path/filepath"
)

func (s *wfwparcio) Connect() error {
	err := s.setSchema()
	if err != nil {
		slog.Error("Cannot get WF/WP schema", "error", err)
	}

	// get Connector to SQLite db
	s.db, err = s.wfwpdb.Connect()
	if err != nil {
		slog.Error("Cannot connect to SQLite database", "error", err)
		return err
	}

	return nil
}

func (s *wfwparcio) setSchema() error {
	defer s.sch.Clean()

	schema, err := s.sch.Fetch()
	if err != nil {
		slog.Error("Cannot fetch schema", "error", err)
		return err
	}

	schFile := filepath.Join(s.cfg.CacheWfwpDir, "schema.sql")
	err = os.WriteFile(schFile, schema, 0644)
	if err != nil {
		slog.Error("Cannot write schema file", "error", err)
		return err
	}
	return nil
}
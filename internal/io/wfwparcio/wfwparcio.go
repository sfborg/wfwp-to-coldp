package wfwparcio

import (
	"database/sql"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"

	"wfwp-to-coldp/internal/ent/wfwparc"
	"wfwp-to-coldp/pkg/config"
	"wfwp-to-coldp/internal/ent/wfwplib"
)

type wfwparcio struct {
	cfg  config.Config
	sch  wfwplib.Schema
	wfwpdb wfwplib.DB
	db   *sql.DB
}

// New creates an instance of WFWPArchive store
func New(cfg config.Config, sch wfwplib.Schema, wfwpdb wfwplib.DB) wfwparc.Archive {
	return &wfwparcio{cfg: cfg, sch: sch, wfwpdb: wfwpdb}
}

func (s *wfwparcio) Exists() bool {
	if s.db == nil {
		return false
	}

	q := "SELECT id FROM version"

	var id string
	err := s.db.QueryRow(q).Scan(&id)
	if err != nil {
		slog.Error("Cannot get data from archive", "error", err)
		return false
	}
	if id == "" {
		slog.Error("Archive version is empty")
		return false
	}

	return true
}

func (s *wfwparcio) Close() error {
	if s.db == nil {
		return nil
	}
	return s.wfwpdb.Close()
}

func (s *wfwparcio) Export(outPath string) error {
	if !s.Exists() {
		return errors.New("cannot find WF/WP archive")
	}

	outPath = trimExtentions(outPath)

	// Determine the desired file extension based on configuration
	ext := ".sql"
	if s.cfg.WithBinOutput {
		ext = ".sqlite"
	}
	outPath += ext

	// Perform the export
	err := s.wfwpdb.Export(outPath, s.cfg.WithBinOutput, s.cfg.WithZipOutput)
	if err != nil {
		return err
	}

	return nil
}

func trimExtentions(outPath string) string {
	hasExt := false
	ext := filepath.Ext(outPath)
	var trimmed string
	if ext == ".zip" {
		hasExt = true
		outPath = strings.TrimSuffix(outPath, ext)
		trimmed += ext
		ext = filepath.Ext(outPath)
	}
	if ext == ".sql" || ext == ".sqlite" {
		hasExt = true
		outPath = strings.TrimSuffix(outPath, ext)
		trimmed = ext + trimmed
	}
	if hasExt {
		slog.Warn("Trimmed extentions from output File", "ext", trimmed)
	}
	return outPath
}
package fwfwp

import (
	wfwpConfig "wfwp-to-coldp/pkg/config"
	"wfwp-to-coldp/internal/ent/wfwp"
	"wfwp-to-coldp/internal/io/arcio"
	"wfwp-to-coldp/internal/ent/wfwparc"
	"wfwp-to-coldp/pkg/config"
)

type fwfwp struct {
	cfg config.Config
	s   wfwparc.Archive
}

func New(cfg config.Config, sfgarc wfwparc.Archive) FromWFWP {
	res := fwfwp{
		cfg: cfg,
		s:   sfgarc,
	}
	return &res
}

// GetWFWP reads a WF/WP Archive from a file, preparing it for ingestion.
func (fw *fwfwp) GetWFWP(path string) (wfwp.Archive, error) {
	opts := []wfwpConfig.Option{
	}
	cfg := wfwpConfig.New(opts...)
	c := arcio.New(cfg, path)
	// Resets cache for wfwp working dir
	err := c.ResetCache()
	if err != nil {
		return nil, err
	}
	err = c.Extract()
	if err != nil {
		return nil, err
	}
	err = c.DirInfo()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ImportCoLDP converts a coldp.Archive to a Species File Group Archive
// database.
func (fw *fwfwp) ImportWFWP(w wfwp.Archive) error {
	err := fw.importData(w)
	if err != nil {
		return err
	}
	return nil
}

// ExportCOLDP writes a COLDP Archive to a file.
func (fc *fwfwp) ExportCOLDP(outputPath string) error {
	err := fc.s.Export(outputPath)
	if err != nil {
		return err
	}

	return nil
}
package arcio

import (
	"wfwp-to-coldp/pkg/config"
	"wfwp-to-coldp/internal/ent/wfwp"
)

type arcio struct {
	// path to the WF/WP archive file
	path string

	// cfg is the configuration
	cfg config.Config

	// metaType can be YAML or JSON
	metaType wfwp.FileType

	// dataPaths contains file paths to dataPaths files.
	// key is low-case type of data, the value is the file path.
	dataPaths map[wfwp.DataType]string
}

func New(cfg config.Config, path string) wfwp.Archive {
	res := arcio{
		path:      path,
		cfg:       cfg,
		dataPaths: make(map[wfwp.DataType]string),
	}

	return &res
}

// DataPaths returns map of low-case names of data files without extensions
// and the path to these files.
func (a *arcio) DataPaths() map[wfwp.DataType]string {
	return a.dataPaths
}

func (a *arcio) Config() config.Config {
	return a.cfg
}
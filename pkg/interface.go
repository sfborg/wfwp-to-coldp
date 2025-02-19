package fwfwp

import (
	"wfwp-to-coldp/internal/ent/wfwp"
)

// FromWFWP provies methods to convert WFWP Archive to Species File Group
// Archive.
type FromWFWP interface {
	// GetWFWP extracts files of WF/WP Archive, preparing it for ingestion.
	GetWFWP(file string) (wfwp.Archive, error)

	// ImportWFWP converts a wfwp.Archive to a Species File Group Archive
	// database.
	ImportWFWP(arc wfwp.Archive) error

	// ExportCOLDP writes a COLDP Archive to a file.
	ExportCOLDP(outputPath string) error
}
package wfwparc

import "wfwp-to-coldp/internal/ent/wfwp"

type Archive interface {
	Exists() bool
	Connect() error
	Close() error

	InsertRows(data []wfwp.Row) error

	Export(outPath string) error
}

type DataWriter interface {
	Write([]DataWriter) error
}

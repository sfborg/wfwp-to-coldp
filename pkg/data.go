package fwfwp

import (
	"log/slog"

	"wfwp-to-coldp/internal/ent/wfwp"
	"wfwp-to-coldp/internal/ent/wfwparc"
)

func (fc *fwfwp) importData(c wfwp.Archive) error {
	var err error
	paths := c.DataPaths()

	if res, ok := paths[wfwp.RowDT]; ok {
		slog.Info("Importing Rows")
		if err = importData(fc, res, c, insertRows); err != nil {
			return err
		}
	}
	return nil
}

func insertRows(s wfwparc.Archive, data []wfwp.Row) error {
	return s.InsertRows(data)
}

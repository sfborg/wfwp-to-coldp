package wfwparcio

import "wfwp-to-coldp/internal/ent/wfwp"

func (s *wfwparcio) InsertRows(data []wfwp.Row) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
	INSERT INTO taxon
		(
		id, rank, number, scientific_name, reference, common_name,
		distribution, synonyms, status, remarks, photo,
		orientation, photographer
		)
	VALUES (
		?,?,?,?,?,?,?,?,?,?,?,?,?,?
		)
`)
	if err != nil {
		return err
	}

	for _, t := range data {
		_, err = stmt.Exec(
			t.Rank, t.Number, t.ScientificName, t.Reference, t.CommonName,
			t.Distribution, t.Synonyms, t.Status, t.Remarks, t.Photo,
			t.Orientation, t.Photographer,
		)
		if err != nil {
			return err
		}
	}

	stmt.Close()
	return tx.Commit()
}
package wfwp

// Row contains information about a World Fern/World Plant taxon
type Row struct {
	// Rank is the taxonomic rank of the name.
	Rank Rank

	// Number is a psudo-identifier for higher classification names?
	Number string

	// Name is a scientific name including the authorship
	ScientificName string

	// Reference is a literature reference for the scientific name
	Reference string

	// CommonName contains common names for the taxon
	CommonName string

	// Distribution contains distribution information for the taxon
	Distribution string

	// Synonyms contains synonyms for the taxon
	Synonyms string

	// Status is the taxonomic status of the name  TODO: confirm?
	Status string

	// Remarks contains additional notes about the taxon  TODO: confirm?
	Remarks string

	// ConservationStatus contains conservation status information  TODO: confirm?
	ConservationStatus string

	// Photo is a link to a photo of the taxon  TODO: confirm?
	Photo string

	// Orientation is the orientation of the photo  TODO: confirm?
	Orientation string

	// Photographer is the name of the photographer  TODO: confirm?
	Photographer string
}

// Load processes a slice of strings into Name object using
// field names from headers.
func (r Row) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	r.Rank = NewRank(row["Taxon"])
	r.Number = row["Number"]
	r.ScientificName = row["Name"]
	r.Reference = row["Literature"]
	r.CommonName = row["TrivialName"]
	r.Distribution = row["Distribution"]
	r.Synonyms = row["Synonyms"]
	r.Status = row["Status"]
	r.Remarks = row["Remarks"]
	r.ConservationStatus = row["ConservationStatus"]
	r.Photo = row["Photo"]
	r.Orientation = row["Orientation"]
	r.Photographer = row["Photographer"]
	return r, warning
}
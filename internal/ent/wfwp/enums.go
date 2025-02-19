package wfwp

import "strings"

// FileType represents the type of file.
type FileType int

// Constants for different file types.
const (
	UnknownFileType FileType = iota
	JSON
	YAML
	CSV
	TSV
	// pipe-separated file
	PSV
	// json-based Citation Style Language
	JSONCSL
	BIBTEX
)

// DataType provides types of data files in CoLDP archive.
// It is used to convert a file name to the type of information
// provided in that file.
type DataType int

// Constants for different data types.
const (
	UnkownDT DataType = iota
	RowDT
)

// FileFormats returns a list of supported file formats for a given DataType.
func (dt DataType) FileFormats() []FileType {
	switch dt {
	case RowDT:
		return []FileType{CSV, PSV, TSV}
	default:
		return []FileType{UnknownFileType}
	}
}

// String returns the string representation of the DataType.
func (dt DataType) String() string {
	return DataTypeToString[dt]
}

// StringToDataType maps strings to DataTypes.
var StringToDataType = map[string]DataType{
	"Row":               RowDT,
}

// DataTypeToString maps DataType to string.
var DataTypeToString = func() map[DataType]string {
	res := make(map[DataType]string)
	for k, v := range StringToDataType {
		res[v] = k
	}
	res[UnkownDT] = "Unknown"
	return res
}()

// LowCaseToDataType provides map of low-case strings to DataType.
var LowCaseToDataType = func() map[string]DataType {
	res := make(map[string]DataType)
	for k, v := range StringToDataType {
		res[strings.ToLower(k)] = v
	}
	return res
}()

// NewDataType creates a new DataType from a string representation.
func NewDataType(s string) DataType {
	s = strings.ToLower(s)
	if dt, ok := LowCaseToDataType[s]; ok {
		return dt
	}
	return UnkownDT
}

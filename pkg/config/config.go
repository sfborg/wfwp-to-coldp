package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/gnfmt"
)

var (
	// jobsNum is the default number of concurrent jobs to run.
	jobsNum = 5
)

// Config is a configuration object for the Catalogue of Life Data
// Package Archive (CoLDP) data processing.
type Config struct {

	// CacheDir is the root path for all cached files.
	CacheDir string

	// CacheWfwpDir is the path WF/WP database.
	CacheWfwpDir string

	// DownloadDir is the path to store downloaded files
	DownloadDir string

	// ExtractDir is the path to store extracted files
	ExtractDir string

	// SchemaPath is the path to the schema file.
	SchemaPath string

	// JobsNum is the number of concurrent jobs to run.
	JobsNum int

	// BatchSize is the number of records to insert in one transaction.
	BatchSize int

	// BadRow sets decision what to do if a row has more/less fields
	// than it should.
	BadRow gnfmt.BadRow

	// WithBinOutput is a flag to output binary SQLite database instead of
	// SQL dump.
	WithBinOutput bool

	// WithZipOutput is a flag to return zipped SFGAarchive outpu.
	WithZipOutput bool

	// WithQuotes tells that coldp file has `"` to escape new lines and
	// delimiters inside fields. If true, RFC-based CSV algorithm is used.
	WithQuotes bool
}

// Option is a function type that allows to standardize how options to
// the configuration are organized.
type Option func(*Config)

// OptCacheDir sets the root path for all temporary files.
func OptCacheDir(s string) Option {
	return func(c *Config) {
		c.CacheDir = s
	}
}

// OptCacheWfwpDir sets the path to store resulting sqlite file with data imported
// from CoLDP file.
func OptCacheWfwpDir(s string) Option {
	return func(c *Config) {
		c.CacheWfwpDir = s
	}
}

// OptJobsNum sets the number of concurrent jobs to run.
func OptJobsNum(n int) Option {
	return func(c *Config) {
		if n < 1 || n > 100 {
			slog.Warn(
				"Unsupported number of jobs (supported: 1-100). Using default value",
				"bad-input", n, "default", jobsNum,
			)
			n = jobsNum
		}
		c.JobsNum = n
	}
}

func OptBadRow(br gnfmt.BadRow) Option {
	return func(c *Config) {
		c.BadRow = br
	}
}

// OptWithBinOutput sets output as binary SQLite file.
func OptWithBinOutput(b bool) Option {
	return func(c *Config) {
		c.WithBinOutput = b
	}
}

// OptWithZipOutput sets output as binary SQLite file.
func OptWithZipOutput(b bool) Option {
	return func(c *Config) {
		c.WithZipOutput = b
	}
}

// OptWithQuotes tells reader that CoLDP file uses quotes in CSV to help
// data integrity when a field contains either new lines, or designated
// field delimiters.
func OptWithQuotes(b bool) Option {
	return func(c *Config) {
		c.WithQuotes = b
	}
}


// New creates a new Config object with default values, and allows to
// override them with options.
func New(opts ...Option) Config {
	tmpDir := os.TempDir()
	path, err := os.UserCacheDir()
	if err != nil {
		path = tmpDir
	}
	path = filepath.Join(path, "wfwp-to-coldp")

	res := Config{
		CacheDir:    path,
		JobsNum:     jobsNum,
		BatchSize:   50_000,
		BadRow:      gnfmt.ErrorBadRow,
	}

	for _, opt := range opts {
		opt(&res)
	}

	res.CacheWfwpDir = filepath.Join(res.CacheDir, "wfwp")
	res.DownloadDir = filepath.Join(res.CacheDir, "download")
	res.ExtractDir = filepath.Join(res.CacheDir, "extract")
	res.SchemaPath = res.CacheDir
	return res
}
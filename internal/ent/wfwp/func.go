package wfwp

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"wfwp-to-coldp/pkg/config"
	"github.com/gnames/gnfmt/gncsv"
	csvConfig "github.com/gnames/gnfmt/gncsv/config"
	"github.com/gnames/gnlib"
)

func Read[T DataLoader](
	cfg config.Config,
	path string,
	chOut chan T) error {
	chIn := make(chan []string)
	var wg sync.WaitGroup
	wg.Add(1)

	opts := []csvConfig.Option{
		csvConfig.OptPath(path),
		csvConfig.OptBadRowMode(cfg.BadRow),
	}
	csvCfg, err := csvConfig.New(opts...)
	headers := gnlib.Map(csvCfg.Headers, func(s string) string {
		return strings.ToLower(s)
	})

	go func() {
		defer wg.Done()
		for row := range chIn {
			var t T
			dl, warn := t.Load(headers, row)
			t = dl.(T)
			if warn != nil {
				slog.Warn("Cannot read row", "warn", warn, "row", row)
			}
			chOut <- t
		}
	}()

	if err != nil {
		return err
	}

	csv := gncsv.New(csvCfg)
	_, err = csv.Read(context.Background(), chIn)
	if err != nil {
		return err
	}
	close(chIn)
	wg.Wait()
	close(chOut)

	return nil
}

type FieldNumberWarning struct {
	FieldsNum, RowFieldsNum int
	Row                     []string
	Message                 string
}

func (w *FieldNumberWarning) Error() string {
	return w.Message
}

func ToInt(s string) int {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	res, _ := strconv.Atoi(s)
	return res
}

func ToBool(s string) bool {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	b := s[0]
	switch b {
	case '1', 'y', 't':
		return true
	default:
		return false
	}
}

func RowToMap(headers, row []string) (map[string]string, error) {
	diff := len(headers) - len(row)
	var warning error
	if diff > 0 {
		for range diff {
			row = append(row, "")
		}
		warning = &FieldNumberWarning{
			FieldsNum:    len(headers),
			RowFieldsNum: len(row),
			Row:          row,
			Message:      fmt.Sprintf("not enough fields, filled with empty strings"),
		}
	}
	if diff < 0 {
		warning = &FieldNumberWarning{
			FieldsNum:    len(headers),
			RowFieldsNum: len(row),
			Row:          row,
			Message:      fmt.Sprintf("too many fields, extras will be ignored"),
		}
	}

	res := make(map[string]string)
	for i := range headers {
		if i < len(row) { // Prevent index out of range
			fld := strings.ToLower(headers[i])
			fld = strings.ReplaceAll(fld, "_", "")
			res[headers[i]] = row[i]
		}
	}

	return res, warning
}
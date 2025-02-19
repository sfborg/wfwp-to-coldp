package arcio

import (
	"wfwp-to-coldp/internal/ent/wfwp"
	"github.com/gnames/gnsys"
)

func (s *arcio) ResetCache() error {
	err := s.emptyCacheDir()
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(s.cfg.DownloadDir)
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(s.cfg.ExtractDir)
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(s.cfg.CacheWfwpDir)
	if err != nil {
		return err
	}

	return nil
}

func (s *arcio) emptyCacheDir() error {
	switch gnsys.GetDirState(s.cfg.CacheDir) {
	case gnsys.DirAbsent:
		return gnsys.MakeDir(s.cfg.CacheDir)
	case gnsys.DirEmpty:
		return nil
	case gnsys.DirNotEmpty:
		return gnsys.CleanDir(s.cfg.CacheDir)
	default:
		return &wfwp.ErrBadDir{Dir: s.cfg.CacheDir}
	}
}
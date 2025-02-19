package fwfwp

import (
	"context"

	"wfwp-to-coldp/internal/ent/wfwp"
	"wfwp-to-coldp/internal/ent/wfwparc"
	"golang.org/x/sync/errgroup"
)

func importData[T wfwp.DataLoader](
	fc *fwfwp,
	path string,
	c wfwp.Archive,
	insertFunc func(wfwparc.Archive, []T) error) error {
	chIn := make(chan T)

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	defer cancel()

	g.Go(func() error {
		return insert(fc.s, fc.cfg.BatchSize, chIn, insertFunc)
	})

	err := wfwp.Read(c.Config(), path, chIn)
	if err != nil {
		return err
	}
	if err = g.Wait(); err != nil {
		return err
	}

	return nil
}

func insert[T wfwp.DataLoader](
	s wfwparc.Archive,
	batchSize int,
	ch <-chan T,
	insertFunc func(wfwparc.Archive, []T) error,
) error {
	var err error
	names := make([]T, 0, batchSize)
	var count int

	for n := range ch {
		count++
		names = append(names, n)
		if count == batchSize {
			err = insertFunc(s, names)
			count = 0
			names = names[:0]
			if err != nil {
				return err
			}
		}
	}

	err = insertFunc(s, names[:count])
	if err != nil {
		return err
	}
	return nil
}
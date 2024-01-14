package entity

import (
	"context"
	fsFactory "github.com/dddplayer/hugoverse/internal/domain/fs/factory"
	"github.com/spf13/afero"
	"os"
)

func newPagesCollector(proc pagesCollectorProcessorProvider, fs afero.Fs) *pagesCollector {
	return &pagesCollector{
		proc: proc,
		fs:   fs,
	}
}

type pagesCollector struct {
	fs   afero.Fs
	proc pagesCollectorProcessorProvider
}

// Collect pages.
func (c *pagesCollector) Collect() (collectErr error) {
	c.proc.Start(context.Background())
	defer func() {
		err := c.proc.Wait()
		if collectErr == nil {
			collectErr = err
		}
	}()

	collectErr = c.collectDir("")
	return
}

func (c *pagesCollector) collectDir(dirname string) error {
	fi, err := c.fs.Stat(dirname)

	if err != nil {
		if os.IsNotExist(err) {
			// May have been deleted.
			return nil
		}
		return err
	}

	//TODO 3
	fsFactory.NewWalkway(c.fs, dirname, func(item any) error {

	})

	return nil
}

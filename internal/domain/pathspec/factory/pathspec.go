package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/factory"
	"github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
)

func NewPathSpec(fs fs.Fs, cfg config.Provider) (*entity.PathSpec, error) {
	p, err := newPaths(fs, cfg)
	if err != nil {
		return nil, err
	}

	bfs, err := factory.NewBaseFS(p)
	if err != nil {
		return nil, err
	}

	ps := &entity.PathSpec{
		Paths:  p,
		BaseFs: bfs,
		Fs:     fs,
		Cfg:    cfg,
	}

	// BasePath not supported yet

	return ps, nil
}

package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
	"path/filepath"
)

func newPaths(fs fs.Fs, cfg config.Provider) (*entity.Paths, error) {
	workingDir := filepath.Clean(cfg.GetString("workingDir"))

	p := &entity.Paths{
		WorkingDir: workingDir,
	}
	if cfg.IsSet("allModules") {
		p.AllModules = cfg.Get("allModules").(config.Modules)
	}

	return p, nil
}
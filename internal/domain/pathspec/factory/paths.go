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
		Fs:         fs,
		WorkingDir: workingDir,
	}
	if cfg.IsSet("modules") {
		p.AllModules = cfg.Get("modules").(config.Modules)
	}

	return p, nil
}

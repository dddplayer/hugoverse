package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
)

func NewFs(path string, cfg config.Provider) fs.Fs {
	cfg.Set("workingDir", path)
	cfg.Set("publishDir", "public")

	return entity.NewFs(cfg.GetString("workingDir"), cfg.GetString("publishDir"))
}

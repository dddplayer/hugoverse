package application

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/config/entity"
	"path"
)

func AllConfigurationInformation(projPath string) (config.Provider, error) {
	c := entity.Config{Path: path.Join(projPath, "config.toml")}

	return c.Load()
}

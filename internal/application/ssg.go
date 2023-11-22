package application

import (
	"github.com/dddplayer/hugoverse/internal/domain/config/entity"
	"github.com/dddplayer/hugoverse/pkg/log"
	"path"
)

func GenerateStaticSite(projPath string, logger log.Logger) error {
	c := entity.Config{
		Path:   path.Join(projPath, "config.toml"),
		Logger: logger,
	}
	provider, err := c.Load()
	if err != nil {
		return err
	}

	logger.Printf("theme: %s", provider.GetString("theme"))
	logger.Printf("modules: %s", provider.Get("modules"))

	return nil
}

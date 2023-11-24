package application

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/config/entity"
	depsFactory "github.com/dddplayer/hugoverse/internal/domain/deps/factory"
	fsFactory "github.com/dddplayer/hugoverse/internal/domain/fs/factory"
	hugoSitesFactory "github.com/dddplayer/hugoverse/internal/domain/hugosites/factory"
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

	fs := fsFactory.NewFs(projPath, provider)
	depsCfg := depsFactory.NewDepsCfg(provider, fs)

	hugoSites, err := hugoSitesFactory.NewHugoSites(depsCfg, logger)
	if err != nil {
		return err
	}
	fmt.Println(hugoSites)

	return nil
}

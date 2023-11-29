package factory

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/deps/entity"
	"github.com/dddplayer/hugoverse/internal/domain/deps/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	hsVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	psFactory "github.com/dddplayer/hugoverse/internal/domain/pathspec/factory"
)

func NewDepsCfg(cfg config.Provider, fs fs.Fs) *valueobject.DepsCfg {
	return &valueobject.DepsCfg{
		Fs:  fs,
		Cfg: cfg,
	}
}

func New(cfg valueobject.DepsCfg) (deps.Deps, error) {
	var (
		originFs = cfg.Fs
	)

	if cfg.TemplateProvider == nil {
		panic("Must have a TemplateProvider")
	}
	if originFs == nil {
		// Default to the production file system.
		panic("Must get originFs ready: deps.New")
	}

	if cfg.MediaTypes == nil {
		cfg.MediaTypes = hsVO.DefaultTypes
	}

	if cfg.OutputFormats == nil {
		cfg.OutputFormats = hsVO.DefaultFormats
	}

	ps, err := psFactory.NewPathSpec(originFs, cfg.Language)
	if err != nil {
		return nil, fmt.Errorf("create PathSpec: %w", err)
	}

	// TODO
	//contentSpec, err := helpers.NewContentSpec(cfg.Language, ps.BaseFs.Content.Fs)
	//if err != nil {
	//	return nil, err
	//}
	//

	d := &entity.Deps{
		Cfg:                 cfg.Language,
		Language:            cfg.Language,
		Site:                cfg.Site, // nil
		OutputFormatsConfig: cfg.OutputFormats,
		TemplateProvider:    cfg.TemplateProvider,

		PathSpec: ps,
	}

	return d, nil
}

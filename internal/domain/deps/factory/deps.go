package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/deps/entity"
	"github.com/dddplayer/hugoverse/internal/domain/deps/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	hugositesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
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
		cfg.MediaTypes = hugositesVO.DefaultTypes
	}

	if cfg.OutputFormats == nil {
		cfg.OutputFormats = hugositesVO.DefaultFormats
	}

	//TODO
	//ps, err := helpers.NewPathSpec(originFs, cfg.Language)
	//if err != nil {
	//	return nil, fmt.Errorf("create PathSpec: %w", err)
	//}
	//
	//resourceSpec, err := resources.NewSpec(ps, cfg.OutputFormats, cfg.MediaTypes)
	//if err != nil {
	//	return nil, err
	//}
	//
	//contentSpec, err := helpers.NewContentSpec(cfg.Language, ps.BaseFs.Content.Fs)
	//if err != nil {
	//	return nil, err
	//}
	//
	//sp := source.NewSourceSpec(ps, nil, originFs.Source)

	d := &entity.Deps{
		Cfg:                 cfg.Language,
		Language:            cfg.Language,
		Site:                cfg.Site,
		OutputFormatsConfig: cfg.OutputFormats,
		TemplateProvider:    cfg.TemplateProvider,
	}

	return d, nil
}

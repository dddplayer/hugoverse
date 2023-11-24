package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/deps/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/template"
)

type Deps struct {
	Cfg              *valueobject.DepsCfg
	TemplateProvider deps.ResourceProvider

	OutputFormatsConfig hugoSitesVO.Formats

	// The templates to use. This will usually implement the full tpl.TemplateManager.
	TemplateHandler template.Handler
}

func NewDeps(cfg config.Provider, fs fs.Fs) *Deps {
	return &Deps{
		Cfg: &valueobject.DepsCfg{
			Fs:  fs,
			Cfg: cfg,
		},
	}
}

func (d *Deps) Config() deps.Cfg {
	return d.Cfg
}

func (d *Deps) Tmpl() template.Handler {
	return d.TemplateHandler
}

package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	lang "github.com/dddplayer/hugoverse/internal/domain/language/entity"
	psEntity "github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
	"github.com/dddplayer/hugoverse/internal/domain/site"
	"github.com/dddplayer/hugoverse/internal/domain/template"
)

type Deps struct {
	Cfg                 config.Provider
	OutputFormatsConfig hugoSitesVO.Formats
	Language            *lang.Language

	// The site building.
	Site site.Site

	TemplateProvider deps.ResourceProvider

	// The templates to use. This will usually implement the full tpl.TemplateManager.
	TemplateHandler template.Handler

	// The PathSpec to use
	PathSpec *psEntity.PathSpec
}

func (d *Deps) Tmpl() template.Handler {
	return d.TemplateHandler
}

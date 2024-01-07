package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/config"
	csEntity "github.com/dddplayer/hugoverse/internal/domain/contentspec/entity"
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
	*psEntity.PathSpec

	// The ContentSpec to use
	*csEntity.ContentSpec `json:"-"`
}

func (d *Deps) Tmpl() template.Handler {
	return d.TemplateHandler
}

func (d *Deps) SetTmpl(tmpl template.Handler) {
	d.TemplateHandler = tmpl
}

// LoadResources loads translations and templates.
func (d *Deps) LoadResources() error {
	if err := d.TemplateProvider.Update(d); err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	return nil
}

func (d *Deps) OutputFormats() hugoSitesVO.Formats {
	return d.OutputFormatsConfig
}

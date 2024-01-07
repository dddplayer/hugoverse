package deps

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	tmplFactory "github.com/dddplayer/hugoverse/internal/domain/template/factory"
	"github.com/spf13/afero"
)

type Deps interface {
	LoadResources() error
	PublishFs() afero.Fs
	LayoutFs() afero.Fs

	Tmpl() template.Handler
	SetTmpl(template.Handler)

	OutputFormats() hugoSitesVO.Formats
}

type Cfg interface {
	Provider() config.Provider
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
	Update(deps Deps) error
	Clone(deps Deps) error
}

// DefaultTemplateProvider is a globally available TemplateProvider.
var DefaultTemplateProvider *TemplateProvider

// TemplateProvider manages templates.
type TemplateProvider struct{}

// Update updates the Hugo Template System in the provided Deps
// with all the additional features, templates & functions.
func (*TemplateProvider) Update(d Deps) error {
	exec, err := tmplFactory.NewTemplateExec(d)
	if err != nil {
		return err
	}
	d.SetTmpl(exec)
	//return tmpl.postTransform()
	return nil
}

func (*TemplateProvider) Clone(d Deps) error {
	panic("not implemented")
	return nil
}

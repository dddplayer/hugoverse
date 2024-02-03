package deps

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/spf13/afero"
)

type Deps interface {
	LoadResources() error
	PublishFs() afero.Fs
	LayoutFs() afero.Fs
	ContentFs() afero.Fs

	Tmpl() template.Handler
	SetTmpl(template.Handler)

	OutputFormats() hugoSitesVO.Formats

	ResolveMarkup(in string) string
	GetContentProvider(name string) contentspec.Provider
}

type Cfg interface {
	Provider() config.Provider
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
	Update(deps Deps) error
	Clone(deps Deps) error
}

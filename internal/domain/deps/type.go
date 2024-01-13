package deps

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/template"
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

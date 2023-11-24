package deps

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/template"
)

type Deps interface {
	Config() Cfg
	Tmpl() template.Handler
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
	return nil
}

func (*TemplateProvider) Clone(d Deps) error {
	panic("not implemented")
	return nil
}

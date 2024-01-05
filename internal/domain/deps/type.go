package deps

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/spf13/afero"
)

type Deps interface {
	LoadResources() error
	Tmpl() template.Handler
	PublishFs() afero.Fs
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
	//TODO
	//tmpl, err := newTemplateExec(d)
	//if err != nil {
	//	return err
	//}
	//return tmpl.postTransform()
	return nil
}

func (*TemplateProvider) Clone(d Deps) error {
	panic("not implemented")
	return nil
}

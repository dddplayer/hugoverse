package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/deps"
)

// DefaultTemplateProvider is a globally available TemplateProvider.
var DefaultTemplateProvider *TemplateProvider

// TemplateProvider manages templates.
type TemplateProvider struct{}

// Update updates the Hugo Template System in the provided Deps
// with all the additional features, templates & functions.
func (*TemplateProvider) Update(d deps.Deps) error {
	exec, err := NewTemplateExec(d)
	if err != nil {
		return err
	}
	d.SetTmpl(exec)
	//return tmpl.postTransform()
	return nil
}

func (*TemplateProvider) Clone(d deps.Deps) error {
	panic("not implemented")
	return nil
}

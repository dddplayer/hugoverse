package entity

import "github.com/dddplayer/hugoverse/internal/domain/template"

// HtmlTemplate is a specialized Template from "text/template" that produces a safe
// HTML document fragment.
type HtmlTemplate struct {
	// We could embed the text/template field, but it's safer not to because
	// we need to keep our version of the name space and the underlying
	// template's in sync.
	Text *TextTemplate

	*nameSpace // common to all associated templates
}

// nameSpace is the data structure shared by all templates in an association.
type nameSpace struct {
	Set map[string]*HtmlTemplate
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type. However, it is legal to overwrite elements of the map. The return
// value is the template, so calls can be chained.
func (t *HtmlTemplate) Funcs(funcMap template.FuncMap) *HtmlTemplate {
	t.Text.Funcs(funcMap)
	return t
}

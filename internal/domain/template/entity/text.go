package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/template"
)

// TextTemplate is the representation of a parsed template. The *parse.Tree
// field is exported only for use by html/template and should be treated
// as unexported by all other clients.
type TextTemplate struct {
	Name string

	*common
}

// common holds the information shared by related templates.
type common struct {
	parseFuncs template.FuncMap
}

// Funcs adds the elements of the argument map to the template's function map.
// It must be called before the template is parsed.
// It panics if a value in the map is not a function with appropriate return
// type or if the name cannot be used syntactically as a function in a template.
// It is legal to overwrite elements of the map. The return value is the template,
// so calls can be chained.
func (t *TextTemplate) Funcs(funcMap template.FuncMap) *TextTemplate {
	addFuncs(t.parseFuncs, funcMap)
	return t
}

// addFuncs adds to values the functions in funcs. It does no checking of the input -
// call addValueFuncs first.
func addFuncs(out, in template.FuncMap) {
	for name, fn := range in {
		out[name] = fn
	}
}

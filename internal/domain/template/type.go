package template

import (
	"context"
	"github.com/dddplayer/hugoverse/internal/domain/template/entity"
	"io"
	"reflect"
)

// Template is the common interface between text/template and html/template.
type Template interface {
	Name() string
	Prepare() (*entity.TextTemplate, error)
}

// Handler finds and executes templates.
type Handler interface {
	Finder
	Execute(t Template, wr io.Writer, data any) error
	ExecuteWithContext(ctx context.Context, t Template, wr io.Writer, data any) error
	HasTemplate(name string) bool
}

// Finder finds templates.
type Finder interface {
	Lookup
}

type Lookup interface {
	Lookup(name string) (Template, bool)
}

// Manager manages the collection of templates.
type Manager interface {
	Handler
	FuncGetter
	AddTemplate(name, tpl string) error
	MarkReady() error
}

// FuncGetter allows to find a template func by name.
type FuncGetter interface {
	GetFunc(name string) (reflect.Value, bool)
}

package entity

import (
	"context"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/spf13/afero"
	"io"
	"reflect"
	"sync"
)

type TemplateExec struct {
	d        deps.Deps
	executor template.Executor
	funcs    map[string]reflect.Value

	*templateHandler
}

type templateHandler struct {
	main      *TemplateNamespace
	readyInit sync.Once
	deps.Deps
}

type TemplateNamespace struct {
	PrototypeText *TextTemplate
	PrototypeHTML *HtmlTemplate

	*TemplateStateMap
}

type TemplateStateMap struct {
	mu        sync.RWMutex
	Templates map[string]*TemplateState
}

type templateType int

type TemplateState struct {
	template.Template

	typ       templateType
	parseInfo ParseInfo

	info     templateInfo
	baseInfo templateInfo // Set when a base template is used.
}

type templateInfo struct {
	name     string
	template string
	isText   bool // HTML or plain text template.

	// Used to create some error context in error situations
	fs afero.Fs

	// The filename relative to the fs above.
	filename string

	// The real filename (if possible). Used for logging.
	realFilename string
}

type ParseInfo struct {
	// Set for shortcode Templates with any {{ .Inner }}
	IsInner bool

	// Set for partials with a return statement.
	HasReturn bool

	// Config extracted from template.
	Config ParseConfig
}

type ParseConfig struct {
	Version int
}

type Executor struct {
	Helper template.ExecHelper
}

// ExecuteWithContext Note: The context is currently not fully implemeted in Hugo. This is a work in progress.
func (t *Executor) ExecuteWithContext(ctx context.Context, p template.Preparer, wr io.Writer, data any) error {
	panic("not implemented")
}

type ExecHelper struct {
	Funcs map[string]reflect.Value
}

var (
	zero             reflect.Value
	contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func (t *ExecHelper) GetFunc(ctx context.Context, tmpl template.Preparer, name string) (fn reflect.Value, firstArg reflect.Value, found bool) {
	if fn, found := t.Funcs[name]; found {
		if fn.Type().NumIn() > 0 {
			first := fn.Type().In(0)
			if first.Implements(contextInterface) {
				// TODO(bep) check if we can void this conversion every time -- and if that matters.
				// The first argument may be context.Context. This is never provided by the end user, but it's used to pass down
				// contextual information, e.g. the top level data context (e.g. Page).
				return fn, reflect.ValueOf(ctx), true
			}
		}

		return fn, zero, true
	}
	return zero, zero, false
}

package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	fsFactory "github.com/dddplayer/hugoverse/internal/domain/fs/factory"
	fsVO "github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
	"sync"
)

type TemplateHandler struct {
	Main *TemplateNamespace
	Deps deps.Deps
}

func (t *TemplateHandler) LoadTemplates() error {
	walker := func(path string, fi fsVO.FileMetaInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}

		name := strings.TrimPrefix(filepath.ToSlash(path), "/")
		if err := t.addTemplateFile(name, path); err != nil {
			return err
		}

		return nil
	}

	return fsFactory.NewWalkway(t.Deps.LayoutFs(), "", walker).Walk()
}

func (t *TemplateHandler) addTemplateFile(name, path string) error {
	//TODO 2

	return nil
}

func (t *TemplateHandler) Lookup(name string) (template.Template, bool) {
	tmpl, found := t.Main.Lookup(name)
	if found {
		return tmpl, true
	}

	return nil, false
}

type TemplateNamespace struct {
	PrototypeText *TextTemplate
	PrototypeHTML *HtmlTemplate

	*TemplateStateMap
}

func (t *TemplateNamespace) Lookup(name string) (template.Template, bool) {
	tmpl, found := t.Templates[name]
	if !found {
		return nil, false
	}

	return tmpl, found
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

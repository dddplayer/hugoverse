package contentspec

type ConverterRegistry interface {
	Get(name string) Provider
}

// Provider creates converters.
type Provider interface {
	New(ctx DocumentContext) (Converter, error)
	Name() string
}

// ProviderProvider creates converter providers.
type ProviderProvider interface {
	New() (Provider, error)
}

// DocumentContext holds contextual information about the document to convert.
type DocumentContext struct {
	Document     any // May be nil. Usually a page.Page
	DocumentID   string
	DocumentName string
	Filename     string
}

// Converter wraps the Convert method that converts some markup into
// another format, e.g. Markdown to HTML.
type Converter interface {
	Convert(ctx RenderContext) (Result, error)
}

// RenderContext holds contextual information about the content to render.
type RenderContext struct {
	// Src is the content to render.
	Src []byte

	// Whether to render TableOfContents.
	RenderTOC bool

	// GerRenderer provides hook renderers on demand.
	//GetRenderer hooks.GetRendererFunc
}

// Result represents the minimum returned from Convert.
type Result interface {
	Bytes() []byte
}

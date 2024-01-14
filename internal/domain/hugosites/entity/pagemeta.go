package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"sync"
)

type pageMeta struct {
	// kind is the discriminator that identifies the different page types
	// in the different page collections. This can, as an example, be used
	// to to filter regular pages, find sections etc.
	// Kind will, for the pages available to the templates, be one of:
	// page, home, section, taxonomy and term.
	// It is of string type to make it easy to reason about in
	// the templates.
	kind string

	// This is a standalone page not part of any page collection. These
	// include sitemap, robotsTXT and similar. It will have no pageOutputs, but
	// a fixed pageOutput.
	standalone bool

	//buildConfig pagemeta.BuildConfig

	// Params contains configuration defined in the params section of page frontmatter.
	params map[string]any

	title     string
	linkTitle string

	summary string

	resourcePath string

	weight int

	markup      string
	contentType string

	// whether the content is in a CJK language.
	isCJKLanguage bool

	layout string

	aliases []string

	description string
	keywords    []string

	//urlPaths pagemeta.URLPath

	//resource.Dates

	// Set if this page is bundled inside another.
	bundled bool

	// A key that maps to translation(s) of this page. This value is fetched
	// from the page front matter.
	translationKey string

	// From front matter.
	configuredOutputFormats valueobject.Formats

	// This is the raw front matter metadata that is going to be assigned to
	// the Resources above.
	resourcesMetadata []map[string]any

	//f source.File

	sections []string

	s *Site

	contentConverterInit sync.Once
	contentConverter     contentspec.Converter
}

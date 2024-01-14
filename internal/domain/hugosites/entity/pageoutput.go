package entity

import "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"

// We create a pageOutput for every output format combination, even if this
// particular page isn't configured to be rendered to that format.
type pageOutput struct {
	// Set if this page isn't configured to be rendered to this format.
	render bool

	f valueobject.Format

	// Only set if render is set.
	// Note that this will be lazily initialized, so only used if actually
	// used in template(s).
	//paginator *pagePaginator

	// These interface provides the functionality that is specific for this
	// output format.
	//pagePerOutputProviders
	//page.ContentProvider
	//page.TableOfContentsProvider
	//page.PageRenderProvider

	// May be nil.
	//cp *pageContentOutput
}

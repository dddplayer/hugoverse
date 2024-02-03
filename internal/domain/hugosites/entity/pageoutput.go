package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
)

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
	hugosites.PagePerOutputProviders
	hugosites.ContentProvider
	//page.TableOfContentsProvider

	// May be nil.
	cp *pageContentOutput
}

func newPageOutput(ps *pageState, pp pagePaths, f valueobject.Format, render bool) *pageOutput {
	var targetPathsProvider targetPathsHolder
	var linksProvider hugosites.ResourceLinksProvider

	ft, found := pp.targetPaths[f.Name]
	if !found {
		// Link to the main output format
		ft = pp.targetPaths[pp.firstOutputFormat.Format.Name]
	}
	targetPathsProvider = ft
	linksProvider = ft

	providers := struct {
		hugosites.ResourceLinksProvider
		hugosites.TargetPather
	}{
		linksProvider,
		targetPathsProvider,
	}

	po := &pageOutput{
		f:                      f,
		PagePerOutputProviders: providers,
		ContentProvider:        nil,
		render:                 render,
	}

	return po
}

func (p *pageOutput) initContentProvider(cp *pageContentOutput) {
	if cp == nil {
		return
	}
	p.ContentProvider = cp
	p.cp = cp
}

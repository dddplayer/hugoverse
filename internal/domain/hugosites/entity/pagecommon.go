package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/pkg/compare"
	"github.com/dddplayer/hugoverse/pkg/lazy"
	"sync"
)

type pageCommon struct {
	s *Site
	m *pageMeta

	//bucket  *pagesMapBucket
	//treeRef *contentTreeRef

	// Lazily initialized dependencies.
	init *lazy.Init

	// All of these represents the common parts of a page.Page
	//page.ChildCareProvider
	hugosites.FileProvider
	hugosites.OutputFormatsProvider
	hugosites.PageMetaProvider
	//hugosites.SitesProvider
	//page.TreeProvider
	//resource.LanguageProvider
	//resource.ResourceMetaProvider
	//resource.ResourceParamsProvider
	//resource.ResourceTypeProvider
	compare.Eqer

	// Describes how paths and URLs for this page and its descendants
	// should look like.
	targetPathDescriptor TargetPathDescriptor

	layoutDescriptor     valueobject.LayoutDescriptor
	layoutDescriptorInit sync.Once

	// The parsed page content.
	pageContent

	// Any bundled resources
	//resources            resource.Resources
	//resourcesInit        sync.Once
	//resourcesPublishInit sync.Once

	//translations    page.Pages
	//allTranslations page.Pages

	// Calculated an cached translation mapping key
	translationKey     string
	translationKeyInit sync.Once

	// Will only be set for bundled pages.
	parent *pageState

	// Set in fast render mode to force render a given page.
	forceRender bool
}

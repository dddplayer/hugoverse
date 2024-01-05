package entity

import "github.com/dddplayer/hugoverse/pkg/radixtree"

type PageMap struct {
	S *Site
	*ContentMap
}

type ContentTree struct {
	Name string
	*radixtree.Tree
}

type ContentTrees []*ContentTree

type ContentMap struct {
	// View of regular pages, sections, and taxonomies.
	PageTrees ContentTrees

	// View of pages, sections, taxonomies, and resources.
	BundleTrees ContentTrees

	// Stores page bundles keyed by its path's directory or the base filename,
	// e.g. "blog/post.md" => "/blog/post", "blog/post/index.md" => "/blog/post"
	// These are the "regular pages" and all of them are bundles.
	Pages *ContentTree

	// Section nodes.
	Sections *ContentTree

	// Resources stored per bundle below a common prefix, e.g. "/blog/post__hb_".
	//Resources *ContentTree
}

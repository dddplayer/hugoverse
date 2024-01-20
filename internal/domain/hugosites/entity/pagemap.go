package entity

import (
	"fmt"
	fsVO "github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/pkg/hugio"
	"github.com/dddplayer/hugoverse/pkg/parser/pageparser"
	"github.com/dddplayer/hugoverse/pkg/radixtree"
	"path"
	"path/filepath"
	"strings"
)

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

func (m *ContentMap) AddFilesBundle(header fsVO.FileMetaInfo) error {
	var (
		meta       = header.Meta()
		bundlePath = m.getBundleDir(meta)

		n = m.newContentNodeFromFi(header)
		b = m.newKeyBuilder()

		section string
	)

	// A regular page. Attach it to its section.
	section, _ = m.getOrCreateSection(n, bundlePath) // /abc/
	b = b.WithSection(section).ForPage(bundlePath).Insert(n)

	return nil
}

func (m *ContentMap) getBundleDir(meta *fsVO.FileMeta) string {
	dir := cleanTreeKey(filepath.Dir(meta.Path))

	switch meta.Classifier {
	case fsVO.ContentClassContent:
		return path.Join(dir, meta.TranslationBaseName)
	default:
		return dir
	}
}

func (m *ContentMap) newContentNodeFromFi(fi fsVO.FileMetaInfo) *contentNode {
	return &contentNode{
		fi:   fi,
		path: strings.TrimPrefix(filepath.ToSlash(fi.Meta().Path), "/"),
	}
}

func (m *ContentMap) newKeyBuilder() *cmInsertKeyBuilder {
	return &cmInsertKeyBuilder{m: m}
}

func (m *ContentMap) getOrCreateSection(n *contentNode, s string) (string, *contentNode) {
	k, b := m.getSection(s)

	k = cleanSectionTreeKey(s[:strings.Index(s[1:], "/")+1])

	b = &contentNode{
		path: n.rootSection(),
	}

	m.Sections.Insert(k, b)

	return k, b
}

func (m *ContentMap) getSection(s string) (string, *contentNode) {
	s = AddTrailingSlash(path.Dir(strings.TrimSuffix(s, "/")))

	v, found := m.Sections.Get(s)
	if found {
		return s, v.(*contentNode)
	}
	return "", nil
}

// AddTrailingSlash adds a trailing Unix styled slash (/) if not already
// there.
func AddTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

func cleanSectionTreeKey(k string) string {
	k = cleanTreeKey(k)
	if k != "/" {
		k += "/"
	}

	return k
}

func cleanTreeKey(k string) string {
	k = "/" + strings.ToLower(strings.Trim(path.Clean(filepath.ToSlash(k)), "./"))
	return k
}

type cmInsertKeyBuilder struct {
	m *ContentMap

	err error

	// Builder state
	tree    *ContentTree
	baseKey string // Section or page key
	key     string
}

func (b *cmInsertKeyBuilder) WithSection(s string) *cmInsertKeyBuilder {
	s = cleanSectionTreeKey(s)
	b.newTopLevel()
	b.tree = b.m.Sections
	b.baseKey = s
	b.key = s
	return b
}

func (b *cmInsertKeyBuilder) newTopLevel() {
	b.key = ""
}

const (
	cmBranchSeparator = "__hb_"
	cmLeafSeparator   = "__hl_"
)

func (b cmInsertKeyBuilder) ForPage(s string) *cmInsertKeyBuilder {
	baseKey := b.baseKey
	b.baseKey = s

	if baseKey != "/" {
		// Don't repeat the section path in the key.
		s = strings.TrimPrefix(s, baseKey)
	}
	s = strings.TrimPrefix(s, "/")

	switch b.tree {
	case b.m.Sections:
		b.tree = b.m.Pages
		b.key = baseKey + cmBranchSeparator + s + cmLeafSeparator
	default:
		panic("invalid state")
	}

	return &b
}

func (b *cmInsertKeyBuilder) Insert(n *contentNode) *cmInsertKeyBuilder {
	if b.err == nil {
		b.tree.Insert(b.Key(), n)
	}
	return b
}

func (b *cmInsertKeyBuilder) Key() string {
	switch b.tree {
	case b.m.Sections:
		return cleanSectionTreeKey(b.key)
	default:
		return cleanTreeKey(b.key)
	}
}

// Assemble

func (m *ContentMap) CreateMissingNodes() error {
	// Create missing home and root sections
	rootSections := make(map[string]any)
	rootSections["/"] = true // not found in both sections and pages

	for sect, _ := range rootSections {
		var sectionPath string
		sect = cleanSectionTreeKey(sect)

		_, found := m.Sections.Get(sect)
		if !found {
			mm := &contentNode{path: sectionPath} // ""
			_ = m.Sections.Insert(sect, mm)       // "/"
		}
	}

	return nil
}

func (m *PageMap) AssemblePages() error {
	var err error
	if err = m.AssembleSections(); err != nil {
		return err
	}

	m.Pages.Walk(func(k string, v any) bool {
		n := v.(*contentNode)
		if n.p == nil {
			return false
		}

		_, parent := m.getSection(k)
		if parent == nil {
			panic(fmt.Sprintf("BUG: parent not set for %q", k))
		}

		n.p, err = m.newPageFromContentNode(n)
		if err != nil {
			return true
		}

		return false
	})

	return err
}

func (m *PageMap) AssembleSections() error {
	m.Sections.Walk(func(k string, v any) bool {
		n := v.(*contentNode)
		sections := m.splitKey(k)

		kind := hugosites.KindSection
		if k == "/" {
			kind = hugosites.KindHome
		}
		if n.fi != nil {
			n.p = m.S.newPage(n, kind, sections...)
		}

		return false
	})

	return nil
}

func (m *PageMap) splitKey(k string) []string {
	if k == "" || k == "/" {
		return nil
	}

	return strings.Split(k, "/")[1:]
}

func (m *PageMap) newPageFromContentNode(n *contentNode) (*pageState, error) {
	if n.fi == nil {
		panic("FileInfo must (currently) be set")
	}

	f, err := newFileInfo(n.fi)
	if err != nil {
		return nil, err
	}

	meta := n.fi.Meta()
	content := func() (hugio.ReadSeekCloser, error) {
		return meta.Open()
	}

	s := m.S
	sections := s.sectionsFromFile(f)
	kind := s.kindFromFileInfoOrSections(f, sections)

	metaProvider := &pageMeta{kind: kind, sections: sections, bundled: false, s: s, f: f}
	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	n.p = ps
	r, err := content()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// .md parseResult
	// TODO: parser works way
	parseResult, err := pageparser.Parse(
		r,
		pageparser.Config{EnableEmoji: false},
	)
	if err != nil {
		return nil, err
	}

	ps.pageContent = pageContent{
		source: rawPageContent{
			parsed:         parseResult,
			posMainContent: -1,
			posSummaryEnd:  -1,
			posBodyStart:   -1,
		},
	}

	if err := ps.mapContent(metaProvider); err != nil {
		return nil, err
	}
	metaProvider.applyDefaultValues()
	ps.init.Add(func() (any, error) {
		// TODO 5 ps init
		// page paths
		// page outputs
		return nil, nil
	})

	return ps, nil
}

package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"path"
	"path/filepath"
	"strings"
)

type targetPathsHolder struct {
	paths valueobject.TargetPaths
	valueobject.OutputFormat
}

func (t targetPathsHolder) TargetPaths() valueobject.TargetPaths {
	return t.paths
}

// TargetPathDescriptor describes how a file path for a given resource
// should look like on the file system. The same descriptor is then later used to
// create both the permalinks and the relative links, paginator URLs etc.
//
// The big motivating behind this is to have only one source of truth for URLs,
// and by that also get rid of most of the fragile string parsing/encoding etc.
type TargetPathDescriptor struct {
	Type valueobject.Format
	Kind string

	Sections []string

	// For regular content pages this is either
	// 1) the Slug, if set,
	// 2) the file base name (TranslationBaseName).
	BaseName string

	// Source directory.
	Dir string

	// Typically a language prefix added to file paths.
	PrefixFilePath string

	// Typically a language prefix added to links.
	PrefixLink string

	// If in multihost mode etc., every link/path needs to be prefixed, even
	// if set in URL.
	ForcePrefix bool

	// URL from front matter if set. Will override any Slug etc.
	URL string

	// The expanded permalink if defined for the section, ready to use.
	ExpandedPermalink string
}

func createTargetPathDescriptor(p hugosites.Page) (TargetPathDescriptor, error) {
	var (
		dir      string
		baseName string
	)

	fmt.Println("555 createTargetPathDescriptor: ", p)
	fmt.Println("555 createTargetPathDescriptor: ", p.File())

	if !p.File().IsZero() {
		dir = p.File().Dir()
		baseName = p.File().TranslationBaseName()
	}

	desc := TargetPathDescriptor{
		Kind:        p.Kind(),
		Sections:    p.SectionsEntries(),
		ForcePrefix: false,
		Dir:         dir,
		URL:         "",
		BaseName:    baseName,
	}

	return desc, nil
}

const slash = "/"

func createTargetPaths(d TargetPathDescriptor) (tp valueobject.TargetPaths) {
	if d.Type.Name == "" {
		panic("CreateTargetPath: missing type")
	}

	if d.URL != "" && !strings.HasPrefix(d.URL, "/") {
		// Treat this as a context relative URL
		d.ForcePrefix = true
	}

	pagePath := slash
	fullSuffix := d.Type.MediaType.FullSuffix()

	var (
		pagePathDir string
		link        string
		linkDir     string
	)

	if d.Kind != valueobject.KindPage && d.URL == "" && len(d.Sections) > 0 {
		pagePath = pjoin(d.Sections...)
	}

	if d.Kind == valueobject.KindPage {
		if d.Dir != "" {
			pagePath = pjoin(pagePath, d.Dir)
		}
		if d.BaseName != "" {
			pagePath = pjoin(pagePath, d.BaseName)
		}

		link = pagePath

		pagePathDir = link
		link = link + slash
		linkDir = pagePathDir

		pagePath = pjoin(pagePath, d.Type.BaseName+fullSuffix)

		if !isHtmlIndex(pagePath) {
			link = pagePath
		}

		if d.PrefixFilePath != "" {
			pagePath = pjoin(d.PrefixFilePath, pagePath)
			pagePathDir = pjoin(d.PrefixFilePath, pagePathDir)
		}

		if d.PrefixLink != "" {
			link = pjoin(d.PrefixLink, link)
			linkDir = pjoin(d.PrefixLink, linkDir)
		}

	} else {

		// No permalink expansion etc. for node type pages (for now)
		base := d.Type.BaseName

		pagePathDir = pagePath
		link = pagePath
		linkDir = pagePathDir

		if base != "" {
			pagePath = path.Join(pagePath, addSuffix(base, fullSuffix))
		} else {
			pagePath = addSuffix(pagePath, fullSuffix)
		}

		if !isHtmlIndex(pagePath) {
			link = pagePath
		} else {
			link += slash
		}

		if d.PrefixFilePath != "" {
			pagePath = pjoin(d.PrefixFilePath, pagePath)
			pagePathDir = pjoin(d.PrefixFilePath, pagePathDir)
		}

		if d.PrefixLink != "" {
			link = pjoin(d.PrefixLink, link)
			linkDir = pjoin(d.PrefixLink, linkDir)
		}
	}

	pagePath = pjoin(slash, pagePath)
	pagePathDir = strings.TrimSuffix(path.Join(slash, pagePathDir), slash)

	hadSlash := strings.HasSuffix(link, slash)
	link = strings.Trim(link, slash)
	if hadSlash {
		link += slash
	}

	if !strings.HasPrefix(link, slash) {
		link = slash + link
	}

	linkDir = strings.TrimSuffix(path.Join(slash, linkDir), slash)

	tp.TargetFilename = filepath.FromSlash(pagePath)
	tp.SubResourceBaseTarget = filepath.FromSlash(pagePathDir)
	tp.SubResourceBaseLink = linkDir
	tp.Link = link
	if tp.Link == "" {
		tp.Link = slash
	}

	return
}

// Like path.Join, but preserves one trailing slash if present.
func pjoin(elem ...string) string {
	hadSlash := strings.HasSuffix(elem[len(elem)-1], slash)
	joined := path.Join(elem...)
	if hadSlash && !strings.HasSuffix(joined, slash) {
		return joined + slash
	}
	return joined
}

func isHtmlIndex(s string) bool {
	return strings.HasSuffix(s, "/index.html")
}

func addSuffix(s, suffix string) string {
	return strings.Trim(s, slash) + suffix
}

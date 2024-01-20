package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	fsVO "github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	langEntity "github.com/dddplayer/hugoverse/internal/domain/language/entity"
	"github.com/dddplayer/hugoverse/pkg/paths"
	"path"
	"strings"
)

type Site struct {
	deps.Deps

	// The owning container. When multiple languages, there will be multiple
	// sites .
	H *HugoSites

	Language *langEntity.Language

	// Output formats defined in site config per Page Kind, or some defaults
	// if not set.
	// Output formats defined in Page front matter will override these.
	OutputFormats map[string]valueobject.Formats

	// All the output formats and media types available for this site.
	// These values will be merged from the Hugo defaults, the site config and,
	// finally, the language settings.
	OutputFormatsConfig valueobject.Formats
	MediaTypesConfig    valueobject.Types

	Publisher hugosites.Publisher

	*PageCollections
}

func (s *Site) process() error {
	if err := s.readAndProcessContent(); err != nil {
		err = fmt.Errorf("readAndProcessContent: %w", err)

		return err
	}

	return nil
}

func (s *Site) readAndProcessContent() error {
	proc := newPagesProcessor(s.H)
	c := newPagesCollector(proc, s.Deps.ContentFs())
	if err := c.Collect(); err != nil {
		return err
	}

	return nil
}

func (s *Site) newPage(n *contentNode, kind string, sections ...string) *pageState {
	p, err := newPageFromMeta(
		n,
		&pageMeta{
			s:        s,
			kind:     kind,
			sections: sections,
		})
	if err != nil {
		panic(err)
	}

	return p
}

func (s *Site) sectionsFromFile(fi hugosites.File) []string {
	dirname := fi.Dir()

	dirname = strings.Trim(dirname, paths.FilePathSeparator)
	if dirname == "" {
		return nil
	}
	parts := strings.Split(dirname, paths.FilePathSeparator)

	if fii, ok := fi.(*fileInfo); ok {
		if len(parts) > 0 && fii.FileInfo().Meta().Classifier == fsVO.ContentClassLeaf {
			// my-section/mybundle/index.md => my-section
			return parts[:len(parts)-1]
		}
	}

	return parts
}

func (s *Site) kindFromFileInfoOrSections(fi *fileInfo, sections []string) string {
	if fi.TranslationBaseName() == "_index" {
		if fi.Dir() == "" {
			return hugosites.KindHome
		}
		return s.kindFromSections(sections)
	}

	return hugosites.KindPage
}

func (s *Site) kindFromSections(sections []string) string {
	if len(sections) == 0 {
		return hugosites.KindHome
	}

	return s.kindFromSectionPath(path.Join(sections...))
}

func (s *Site) kindFromSectionPath(sectionPath string) string {
	return hugosites.KindSection
}

func (s *Site) initRenderFormats() {
	// TODO 7
}

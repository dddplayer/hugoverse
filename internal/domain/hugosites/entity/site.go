package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	langEntity "github.com/dddplayer/hugoverse/internal/domain/language/entity"
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

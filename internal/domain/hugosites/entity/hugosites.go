package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/pkg/para"
	"sync"
)

type HugoSites struct {
	deps.Deps

	Sites      []*Site
	NumWorkers int

	contentInit sync.Once
	content     *pageMaps

	// Render output formats for all sites.
	RenderFormats valueobject.Formats
}

func (h *HugoSites) Build() error {
	// process file system to create content map
	err := h.process()
	if err != nil {
		return err
	}

	err = h.assemble()
	if err != nil {
		return err
	}

	err = h.render()
	if err != nil {
		return err
	}

	return nil
}

func (h *HugoSites) process() error {
	firstSite := h.Sites[0]
	return firstSite.process()
}

func (h *HugoSites) assemble() error {
	if err := h.getContentMaps().AssemblePages(); err != nil {
		return err
	}
	return nil
}

func (h *HugoSites) render() error {
	h.withSite(func(s *Site) error {
		s.initRenderFormats()
		return nil
	})

	h.RenderFormats = valueobject.Formats{}
	for _, s := range h.Sites {
		// TODO: remove duplication from different sites
		h.RenderFormats = append(h.RenderFormats, s.RenderFormats...)
	}

	for _, s := range h.Sites {
		// Get page output ready
		if err := s.preparePagesForRender(); err != nil {
			return err
		}
		if err := s.render(); err != nil {
			return err
		}
	}

	return nil
}

func (h *HugoSites) withSite(fn func(s *Site) error) {
	for _, s := range h.Sites {
		if err := fn(s); err != nil {
			panic(err)
		}
	}
}

func (h *HugoSites) getContentMaps() *pageMaps {
	h.contentInit.Do(func() {
		h.content = newPageMaps(h)
	})
	return h.content
}

func newPageMaps(h *HugoSites) *pageMaps {
	mps := make([]*PageMap, len(h.Sites))
	for i, s := range h.Sites {
		mps[i] = s.PageMap
	}
	return &pageMaps{
		workers: para.New(h.NumWorkers),
		pmaps:   mps,
	}
}

package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/deps"
)

type HugoSites struct {
	deps.Deps

	Sites      []*Site
	Init       *HugoSitesInit
	NumWorkers int
}

func (h *HugoSites) Build() error {

	// process file system to create content map
	err := h.process()
	if err != nil {
		return err
	}

	return nil
}

func (h *HugoSites) process() error {
	firstSite := h.Sites[0]
	return firstSite.process()
}

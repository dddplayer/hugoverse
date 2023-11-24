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

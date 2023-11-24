package factory

import (
	depsVO "github.com/dddplayer/hugoverse/internal/domain/deps/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/entity"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
)

// newSite creates a new site with the given configuration.
func newSite(cfg depsVO.DepsCfg) (*entity.Site, error) {
	mediaTypes := valueobject.DecodeTypes()
	formats := valueobject.DecodeFormats(mediaTypes)
	outputFormats := valueobject.CreateSiteOutputFormats(formats)

	s := &entity.Site{
		Language: cfg.Language,

		OutputFormats:       outputFormats,
		OutputFormatsConfig: formats,
		MediaTypesConfig:    mediaTypes,
	}

	return s, nil
}

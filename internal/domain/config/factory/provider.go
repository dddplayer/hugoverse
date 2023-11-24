package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/config/valueobject"
)

func New() config.Provider {
	return &valueobject.DefaultConfigProvider{
		Root: make(config.Params),
	}
}

// NewCompositeConfig creates a new composite Provider with a read-only base
// and a writeable layer.
func NewCompositeConfig(base, layer config.Provider) config.Provider {
	return &valueobject.CompositeConfig{
		Base:  base,
		Layer: layer,
	}
}

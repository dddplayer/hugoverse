package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/config/valueobject"
)

func New() config.Provider {
	return &valueobject.DefaultConfigProvider{
		Root: make(valueobject.Params),
	}
}

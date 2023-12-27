package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
)

type ContentSpec struct {
	Converters contentspec.ConverterRegistry
	Cfg        config.Provider
}

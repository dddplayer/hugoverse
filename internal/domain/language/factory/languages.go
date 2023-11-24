package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/language/entity"
)

func GetLanguages(cfg config.Provider) entity.Languages {
	return entity.Languages{entity.NewDefaultLanguage(cfg)}
}

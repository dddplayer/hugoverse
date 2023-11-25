package valueobject

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	hugoSitesVO "github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	lang "github.com/dddplayer/hugoverse/internal/domain/language/entity"
	"github.com/dddplayer/hugoverse/internal/domain/site"
)

// DepsCfg contains configuration options that can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {
	// The file systems to use
	Fs fs.Fs

	// The configuration to use.
	Cfg config.Provider

	// The Site in use
	Site site.Site

	// The language to use.
	Language *lang.Language

	// The media types configured.
	MediaTypes hugoSitesVO.Types

	// The output formats configured.
	OutputFormats hugoSitesVO.Formats

	// Template handling.
	TemplateProvider deps.ResourceProvider
}

func (dc *DepsCfg) Provider() config.Provider {
	return dc.Cfg
}

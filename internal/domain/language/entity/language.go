package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	configFactory "github.com/dddplayer/hugoverse/internal/domain/config/factory"
)

// Language manages specific-language configuration.
type Language struct {
	Lang   string
	Weight int // for sort

	// If set per language, this tells Hugo that all content files without any
	// language indicator (e.g. my-page.en.md) is in this language.
	// This is usually a path relative to the working dir, but it can be an
	// absolute directory reference. It is what we get.
	// For internal use.
	ContentDir string

	// Global config.
	// For internal use.
	Cfg config.Provider

	// Language specific config.
	// For internal use.
	LocalCfg config.Provider

	// Composite config.
	// For internal use.
	config.Provider

	// Error during initialization. Will fail the buld.
	initErr error
}

// NewDefaultLanguage creates the default language for a config.Provider.
// If not otherwise specified the default is "en".
func NewDefaultLanguage(cfg config.Provider) *Language {
	defaultLang := cfg.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return NewLanguage(defaultLang, cfg)
}

// NewLanguage creates a new language.
func NewLanguage(lang string, cfg config.Provider) *Language {
	localCfg := configFactory.New()
	compositeConfig := configFactory.NewCompositeConfig(cfg, localCfg)

	l := &Language{
		Lang:       lang,
		ContentDir: cfg.GetString("contentDir"),
		Cfg:        cfg,
		LocalCfg:   localCfg,
		Provider:   compositeConfig,
	}

	return l
}

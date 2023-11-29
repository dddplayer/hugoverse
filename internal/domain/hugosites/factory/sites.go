package factory

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/deps"
	depsFactory "github.com/dddplayer/hugoverse/internal/domain/deps/factory"
	depsVO "github.com/dddplayer/hugoverse/internal/domain/deps/valueobject"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/entity"
	langsFactory "github.com/dddplayer/hugoverse/internal/domain/language/factory"
	"github.com/dddplayer/hugoverse/internal/domain/template"
	"github.com/dddplayer/hugoverse/pkg/lazy"
	"github.com/dddplayer/hugoverse/pkg/log"
)

func NewHugoSites(cfg *depsVO.DepsCfg, logger log.Logger) (*entity.HugoSites, error) {
	sites, err := createSitesFromConfig(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("from config: %w", err)
	}

	return newHugoSites(cfg, logger, sites...)
}

func createSitesFromConfig(cfg *depsVO.DepsCfg, log log.Logger) ([]*entity.Site, error) {
	log.Printf("createSitesFromConfig: %s", "start")
	var sites []*entity.Site

	// [en]
	languages := langsFactory.GetLanguages(cfg.Cfg)
	for _, lang := range languages {
		var s *entity.Site
		var err error
		cfg.Language = lang
		log.Printf("newSite: %s", "create site with DepsCfg with language setup")

		s, err = newSite(*cfg)
		if err != nil {
			return nil, err
		}

		sites = append(sites, s)
	}

	log.Printf("createSitesFromConfig: %s", "end")
	return sites, nil
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(cfg *depsVO.DepsCfg, log log.Logger, sites ...*entity.Site) (*entity.HugoSites, error) {
	var initErr error

	log.Printf("newHugoSites: %s", "init HugoSites")
	h := &entity.HugoSites{
		Sites:      sites,
		NumWorkers: 1,
		Init: &entity.HugoSitesInit{
			Layouts: lazy.New(),
		},
	}

	log.Printf("newHugoSites: %s", "add layouts to h.init")
	h.Init.Layouts.Add(func() (any, error) {
		log.Printf("newHugoSites: %s", "h.init register s.Tmpl().MarkReady when deps not set for all sites")
		for _, s := range h.Sites {
			if err := s.Tmpl().(template.Manager).MarkReady(); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	for _, s := range sites {
		s.H = h
	}

	log.Printf("newHugoSites: %s", "configLoader applyDeps")
	if err := applyDeps(cfg, log, sites...); err != nil {
		initErr = fmt.Errorf("add site dependencies: %w", err)
	}

	h.Deps = sites[0].Deps
	if h.Deps == nil {
		return nil, initErr
	}

	return h, initErr

}

func applyDeps(cfg *depsVO.DepsCfg, log log.Logger, sites ...*entity.Site) error {
	log.Printf("applyDeps: %s", "set cfg.TemplateProvider with DefaultTemplateProvider")

	if cfg.TemplateProvider == nil {
		cfg.TemplateProvider = deps.DefaultTemplateProvider
	}

	var (
		d deps.Deps
	)

	for _, s := range sites {
		if s.Deps != nil {
			continue
		}

		onCreated := func(d deps.Deps) error {
			s.Deps = d

			log.Printf("applyDeps-onCreate: %s", "set site publisher as DestinationPublisher")
			//TODO

			// d.Site = s.Info

			return nil
		}

		cfg.Language = s.Language
		cfg.MediaTypes = s.MediaTypesConfig
		cfg.OutputFormats = s.OutputFormatsConfig

		log.Printf("applyDeps: %s", "new deps")

		var err error
		d, err = depsFactory.New(*cfg)
		if err != nil {
			return fmt.Errorf("create deps: %w", err)
		}

		if err := onCreated(d); err != nil {
			return fmt.Errorf("on created: %w", err)
		}

		log.Printf("applyDeps: %s", "deps LoadResources to update template provider, need to make template ready")

		//TODO
		//if err = d.LoadResources(); err != nil {
		//	return fmt.Errorf("load resources: %w", err)
		//}
	}

	return nil
}

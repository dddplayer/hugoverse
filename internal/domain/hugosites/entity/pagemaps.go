package entity

import (
	"context"
	"github.com/dddplayer/hugoverse/pkg/para"
)

type pageMaps struct {
	workers *para.Workers
	pmaps   []*PageMap
}

func (m *pageMaps) AssemblePages() error {
	return m.withMaps(func(pm *PageMap) error {
		if err := pm.CreateMissingNodes(); err != nil {
			return err
		}

		if err := pm.AssemblePages(); err != nil {
			return err
		}

		// Handle any new sections created in the step above.
		if err := pm.AssembleSections(); err != nil {
			return err
		}
		return nil
	})
}

func (m *pageMaps) withMaps(fn func(pm *PageMap) error) error {
	g, _ := m.workers.Start(context.Background())
	for _, pmap := range m.pmaps {
		pm := pmap
		g.Run(func() error {
			return fn(pm)
		})
	}
	return g.Wait()
}

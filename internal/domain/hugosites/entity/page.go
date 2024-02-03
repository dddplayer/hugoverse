package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
	"github.com/dddplayer/hugoverse/pkg/lazy"
)

var (
	nopPageOutput = &pageOutput{
		// TODO, simplify
	}
)

func newPageFromMeta(n *contentNode, metaProvider *pageMeta) (*pageState, error) {
	if metaProvider.f == nil {
		metaProvider.f = NewZeroFile()
	}

	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	metaProvider.setMetadata()
	metaProvider.applyDefaultValues()

	ps.init.Add(func() (any, error) {
		fmt.Println("init 222--- ")

		pp, err := newPagePaths(metaProvider.s, ps, metaProvider)
		if err != nil {
			return nil, err
		}

		makeOut := func(f valueobject.Format, render bool) *pageOutput {
			return newPageOutput(ps, pp, f, render)
		}

		outputFormatsForPage := ps.m.outputFormats()
		// Prepare output formats for all sites.
		// We do this even if this page does not get rendered on
		// its own. It may be referenced via .Site.GetPage and
		// it will then need an output format.
		ps.pageOutputs = make([]*pageOutput, len(ps.s.H.RenderFormats))
		created := make(map[string]*pageOutput)
		for i, f := range ps.s.H.RenderFormats {
			po, found := created[f.Name]
			if !found {
				_, render := outputFormatsForPage.GetByName(f.Name)
				po = makeOut(f, render)
				created[f.Name] = po
			}
			ps.pageOutputs[i] = po
		}
		if err := ps.initCommonProviders(pp); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return ps, err
}

func newPageBase(metaProvider *pageMeta) (*pageState, error) {
	if metaProvider.s == nil {
		panic("must provide a Site")
	}

	s := metaProvider.s

	ps := &pageState{
		pageOutput: nopPageOutput,
		pageCommon: &pageCommon{
			// Simplify:  FileProvider...
			FileProvider:     metaProvider,
			PageMetaProvider: metaProvider,
			init:             lazy.New(),
			m:                metaProvider,
			s:                s,
		},
	}

	return ps, nil
}

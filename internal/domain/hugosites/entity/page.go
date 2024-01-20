package entity

import "github.com/dddplayer/hugoverse/pkg/lazy"

var (
	nopPageOutput = &pageOutput{
		// TODO, simplify
	}
)

func newPageFromMeta(n *contentNode, metaProvider *pageMeta) (*pageState, error) {
	ps, err := newPageBase(metaProvider)
	if err != nil {
		return nil, err
	}

	metaProvider.setMetadata()
	metaProvider.applyDefaultValues()

	ps.init.Add(func() (any, error) {
		// TODO 4: page state init
		// new page paths
		// new page output

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
			init: lazy.New(),
			m:    metaProvider,
			s:    s,
		},
	}

	return ps, nil
}

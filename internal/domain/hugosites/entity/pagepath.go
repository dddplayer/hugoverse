package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/dddplayer/hugoverse/internal/domain/hugosites/valueobject"
)

func newPagePaths(s *Site, p hugosites.Page, pm *pageMeta) (pagePaths, error) {
	fmt.Println("newPagePaths", p)

	targetPathDescriptor, err := createTargetPathDescriptor(p)
	if err != nil {
		return pagePaths{}, err
	}

	outputFormats := pm.outputFormats()
	if len(outputFormats) == 0 {
		fmt.Println("outputFormats is null", outputFormats)

		return pagePaths{}, nil
	}

	pageOutputFormats := make(valueobject.OutputFormats, len(outputFormats))
	targets := make(map[string]targetPathsHolder)

	for i, f := range outputFormats {
		desc := targetPathDescriptor
		desc.Type = f
		paths := createTargetPaths(desc)

		var relPermalink, permalink string

		// If a page is headless or bundled in another,
		// it will not get published on its own and it will have no links.
		// We also check the build options if it's set to not render or have
		// a link.
		if !pm.noLink() && !pm.bundled {
			relPermalink = paths.RelPermalink()
			permalink = paths.PermalinkForOutputFormat()
		}

		pageOutputFormats[i] = valueobject.NewOutputFormat(relPermalink, permalink, f)

		// Use the main format for permalinks, usually HTML.
		permalinksIndex := 0
		targets[f.Name] = targetPathsHolder{
			paths:        paths,
			OutputFormat: pageOutputFormats[permalinksIndex],
		}

	}

	var out valueobject.OutputFormats
	if !pm.noLink() {
		out = pageOutputFormats
	}

	return pagePaths{
		outputFormats:        out,
		firstOutputFormat:    pageOutputFormats[0],
		targetPaths:          targets,
		targetPathDescriptor: targetPathDescriptor,
	}, nil
}

type pagePaths struct {
	outputFormats     valueobject.OutputFormats
	firstOutputFormat valueobject.OutputFormat

	targetPaths          map[string]targetPathsHolder
	targetPathDescriptor TargetPathDescriptor
}

func (l pagePaths) OutputFormats() valueobject.OutputFormats {
	return l.outputFormats
}

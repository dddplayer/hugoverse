package entity

import (
	"errors"
	"github.com/dddplayer/hugoverse/pkg/parser/pageparser"
)

type pageState struct {
	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// This will be shifted out when we start to render a new output format.
	*pageOutput

	// Common for all output formats.
	*pageCommon
}

func (p *pageState) mapContent(meta *pageMeta) error {
	p.cmap = &pageContentMap{
		items: make([]any, 0, 20),
	}
	return p.mapContentForResult(
		p.source.parsed,
		p.cmap,
		meta.markup,
	)
}

func (p *pageState) mapContentForResult(result pageparser.Result, rn *pageContentMap, markup string) error {
	iter := result.Iterator()
	fail := func(err error, i pageparser.Item) error {
		return errors.New("fail fail fail")
	}

Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.IsFrontMatter():
			panic("not implemented front matter yet")
		case it.Type == pageparser.TypeLeadSummaryDivider:
			panic("not implemented lead summary divider yet")
		case it.Type == pageparser.TypeEmoji:
			panic("not implemented emoji yet")
		case it.IsEOF():
			break Loop
		case it.IsError():
			err := fail(errors.New(it.ValStr(result.Input())), it)
			return err

		default:
			rn.AddBytes(it)
		}
	}

	return nil
}

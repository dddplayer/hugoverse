package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
)

// MDProvider is the package entry point.
var MDProvider contentspec.ProviderProvider = provide{}

type provide struct {
	name string
}

func (p provide) New() (contentspec.Provider, error) {
	//TODO, implement me with dddplayer/markdown
	// md := newMarkdown()

	return ConverterProvider{
		name: "markdown",
		create: func(ctx contentspec.DocumentContext) (contentspec.Converter, error) {
			return &mdConverter{}, nil
		},
	}, nil
}

type mdConverter struct {
	//md dddplayer.markdown
}

func (c *mdConverter) Convert(ctx contentspec.RenderContext) (result contentspec.Result, err error) {
	return nil, nil
}

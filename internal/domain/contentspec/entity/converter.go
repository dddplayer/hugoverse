package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
	"strings"
)

type ConverterRegistry struct {
	// Maps name (md, markdown, goldmark etc.) to a converter provider.
	// Note that this is also used for aliasing, so the same converter
	// may be registered multiple times.
	// All names are lower case.
	Converters map[string]contentspec.Provider
}

func (r *ConverterRegistry) Get(name string) contentspec.Provider {
	return r.Converters[strings.ToLower(name)]
}

type ConverterProvider struct {
	name   string
	create func(ctx contentspec.DocumentContext) (contentspec.Converter, error)
}

func (n ConverterProvider) New(ctx contentspec.DocumentContext) (contentspec.Converter, error) {
	return n.create(ctx)
}

func (n ConverterProvider) Name() string {
	return n.name
}

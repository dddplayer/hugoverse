package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
	"github.com/dddplayer/hugoverse/internal/domain/contentspec/entity"
)

// NewContentSpec returns a ContentSpec initialized
// with the appropriate fields from the given config.Provider.
func NewContentSpec(cfg config.Provider) (*entity.ContentSpec, error) {
	spec := &entity.ContentSpec{
		Cfg: cfg,
	}

	// markdown
	converterRegistry, err := NewConverterRegistry()
	if err != nil {
		return nil, err
	}

	spec.Converters = converterRegistry

	return spec, nil
}

func NewConverterRegistry() (contentspec.ConverterRegistry, error) {
	converters := make(map[string]contentspec.Provider)

	add := func(p contentspec.ProviderProvider) error {
		c, err := p.New()
		if err != nil {
			return err
		}

		name := c.Name()
		converters[name] = c

		return nil
	}

	// default
	if err := add(entity.MDProvider); err != nil {
		return nil, err
	}

	return &entity.ConverterRegistry{
		Converters: converters,
	}, nil
}

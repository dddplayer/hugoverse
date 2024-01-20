package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/contentspec"
	"strings"
)

type ContentSpec struct {
	Converters contentspec.ConverterRegistry
	Cfg        config.Provider
}

func (c *ContentSpec) ResolveMarkup(in string) string {
	in = strings.ToLower(in)
	switch in {
	case "md", "markdown", "mdown":
		return "markdown"
	case "html", "htm":
		return "html"
	default:
		if conv := c.Converters.Get(in); conv != nil {
			return conv.Name()
		}
	}
	return ""
}

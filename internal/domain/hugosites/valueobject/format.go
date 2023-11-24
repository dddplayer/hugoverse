package valueobject

import (
	"fmt"
	"strings"
)

// Format represents an output representation, usually to a file on disk.
type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string `json:"name"`

	MediaType Type `json:"-"`

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string `json:"baseName"`
}

// Formats is a slice of Format.
type Formats []Format

// GetByName gets a format by its identifier name.
func (formats Formats) GetByName(
	name string) (f Format, found bool) {
	for _, ff := range formats {
		if strings.EqualFold(name, ff.Name) {
			f = ff
			found = true
			return
		}
	}
	return
}

// HTMLFormat An ordered list of built-in output formats.
var HTMLFormat = Format{
	Name:      "HTML",
	MediaType: HTMLType,
	BaseName:  "index",
}

// DefaultFormats contains the default output formats supported by Hugo.
var DefaultFormats = Formats{
	HTMLFormat,
}

// DecodeFormats takes a list of output format configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeFormats(mediaTypes Types) Formats {
	// Format could be modified by mediaTypes configuration
	// just make it simple for example
	fmt.Println(mediaTypes)

	f := make(Formats, len(DefaultFormats))
	copy(f, DefaultFormats)

	return f
}

func CreateSiteOutputFormats(allFormats Formats) map[string]Formats {
	defaultOutputFormats :=
		createDefaultOutputFormats(allFormats)
	return defaultOutputFormats
}

const (
	KindPage = "page"
	kind404  = "404"
)

func createDefaultOutputFormats(
	allFormats Formats) map[string]Formats {
	htmlOut, _ := allFormats.GetByName(HTMLFormat.Name)

	m := map[string]Formats{
		KindPage: {htmlOut},
		kind404:  {htmlOut},
	}

	return m
}

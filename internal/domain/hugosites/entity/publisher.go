package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/hugosites"
	"github.com/spf13/afero"
)

// DestinationPublisher is the default and currently only publisher in Hugo. This
// publisher prepares and publishes an item to the defined destination, e.g. /public.
type DestinationPublisher struct {
	Fs afero.Fs
	//min minifiers.Client
}

// Publish applies any relevant transformations and writes the file
// to its destination, e.g. /public.
func (p *DestinationPublisher) Publish(d hugosites.Descriptor) error {
	//TODO
	return nil
}

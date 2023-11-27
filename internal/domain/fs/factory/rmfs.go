package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/pkg/radixtree"
)

// NewRootMappingFs creates a new RootMappingFs
// on top of the provided with root mappings with
// some optional metadata about the root.
// Note that From represents a virtual root
// that maps to the actual filename in To.
func newRootMappingFs(rms ...valueobject.RootMapping) *entity.RootMappingFs {
	t := radixtree.New()
	var virtualRoots []valueobject.RootMapping

	for _, rm := range rms {
		key := fs.FilepathSeparator + rm.From
		mappings := valueobject.GetRms(t, key)
		mappings = append(mappings, rm)
		t.Insert(key, mappings)

		virtualRoots = append(virtualRoots, rm)
	}

	t.Insert(fs.FilepathSeparator, virtualRoots)

	return &entity.RootMappingFs{
		RootMapToReal: t,
	}
}

package valueobject

import "github.com/dddplayer/hugoverse/pkg/radixtree"

// RootMapping describes a virtual file or directory mount.
type RootMapping struct {
	// The virtual mount.
	From string
	// The source directory or file.
	To string
	// The base of To. May be empty if an
	// absolute path was provided.
	ToBasedir string
	// Whether this is a mount in the main project.
	IsProject bool
	// The virtual mount point, e.g. "blog".
	path string

	Meta *FileMeta // File metadata (lang etc.)
}

func GetRms(t *radixtree.Tree, key string) []RootMapping {
	var mappings []RootMapping
	v, found := t.Get(key)
	if found {
		mappings = v.([]RootMapping)
	}
	return mappings
}

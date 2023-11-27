package entity

import (
	dfs "github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/dddplayer/hugoverse/pkg/radixtree"
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"path"
)

// A RootMappingFs maps several roots into one.
// Note that the root of this filesystem
// is directories only, and they will be returned
// in Readdir and Readdirnames
// in the order given.
type RootMappingFs struct {
	fs            afero.Fs
	RootMapToReal *radixtree.Tree
}

func (m *RootMappingFs) Abs(name string) []string {
	mappings := valueobject.GetRms(m.RootMapToReal, name)

	var paths []string
	for _, m := range mappings {
		paths = append(paths, path.Join(
			m.ToBasedir, m.To))
	}
	return paths
}

func (m *RootMappingFs) Fss() []fs.StatFS {
	mappings := valueobject.GetRms(m.RootMapToReal, dfs.FilepathSeparator)

	var fss []fs.StatFS
	for _, m := range mappings {
		fss = append(fss, os.DirFS(
			path.Join(m.ToBasedir, m.To)).(fs.StatFS))
	}
	return fss
}

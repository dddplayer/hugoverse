package entity

import (
	"fmt"
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
	Fs            afero.Fs
	RootMapToReal *radixtree.Tree
}

func (m *RootMappingFs) Dirs(base string) ([]valueobject.FileMetaInfo, error) {
	base = dfs.FilepathSeparator + base
	roots := m.getRootsWithPrefix(base)

	if roots == nil {
		return nil, nil
	}

	fss := make([]valueobject.FileMetaInfo, len(roots))
	for i, r := range roots {
		bfs := afero.NewBasePathFs(m.Fs, r.To)
		bfs = decorateDirs(bfs, r.Meta)

		fi, err := bfs.Stat("")
		if err != nil {
			return nil, fmt.Errorf("RootMappingFs.Dirs: %w", err)
		}

		fss[i] = fi.(valueobject.FileMetaInfo)
	}

	return fss, nil
}

func (m *RootMappingFs) getRootsWithPrefix(prefix string) []valueobject.RootMapping {
	return valueobject.GetRms(m.RootMapToReal, prefix) // /content
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

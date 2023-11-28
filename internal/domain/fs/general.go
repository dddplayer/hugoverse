package fs

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

var FilepathSeparator = string(filepath.Separator)

const (
	ComponentFolderArchetypes = "archetypes"
	ComponentFolderStatic     = "static"
	ComponentFolderLayouts    = "layouts"
	ComponentFolderContent    = "content"
	ComponentFolderData       = "data"
	ComponentFolderAssets     = "assets"
	ComponentFolderI18n       = "i18n"
)

var (
	ComponentFolders = []string{
		ComponentFolderArchetypes,
		ComponentFolderStatic,
		ComponentFolderLayouts,
		ComponentFolderContent,
		ComponentFolderData,
		ComponentFolderAssets,
		ComponentFolderI18n,
	}
)

// LstatIfPossible if the filesystem supports it, use Lstat, else use fs.Stat
func LstatIfPossible(fs afero.Fs, path string) (os.FileInfo, bool, error) {
	if lfs, ok := fs.(afero.Lstater); ok {
		fi, b, err := lfs.LstatIfPossible(path)
		return fi, b, err
	}
	fi, err := fs.Stat(path)
	return fi, false, err
}

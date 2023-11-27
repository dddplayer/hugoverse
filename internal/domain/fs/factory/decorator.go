package factory

import (
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

// newBaseFileDecorator decorates the given Fs to provide the real filename
// and an Opener func.
func newBaseFileDecorator(originFs afero.Fs) afero.Fs {
	ffs := &entity.BaseFileDecoratorFs{Fs: originFs}

	decorator := func(fi os.FileInfo, filename string) (os.FileInfo, error) {
		// Store away the original in case it's a symlink.
		meta := valueobject.NewFileMeta()
		meta.Name = fi.Name()

		if fi.IsDir() {
			meta.JoinStatFunc = func(name string) (valueobject.FileMetaInfo, error) {
				joinedFilename := filepath.Join(filename, name)
				fi, _, err := fs.LstatIfPossible(originFs, joinedFilename)
				if err != nil {
					return nil, err
				}

				fi, err = ffs.Decorate(fi, joinedFilename)
				if err != nil {
					return nil, err
				}

				return fi.(valueobject.FileMetaInfo), nil
			}
		}

		opener := func() (afero.File, error) {
			return ffs.Open(filename)
		}

		fim := valueobject.DecorateFileInfo(fi, ffs, opener, filename, "", meta)

		return fim, nil
	}

	ffs.Decorate = decorator
	return ffs
}

package valueobject

import (
	"github.com/dddplayer/hugoverse/pkg/htime"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"time"
)

func NewDirNameOnlyFI(name string, modTime time.Time) *DirNameOnlyFileInfo {
	return &DirNameOnlyFileInfo{name: name, modTime: modTime}
}

type DirNameOnlyFileInfo struct {
	name    string
	modTime time.Time
}

func (fi *DirNameOnlyFileInfo) Name() string {
	return fi.name
}

func (fi *DirNameOnlyFileInfo) Size() int64 {
	panic("not implemented")
}

func (fi *DirNameOnlyFileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (fi *DirNameOnlyFileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi *DirNameOnlyFileInfo) IsDir() bool {
	return true
}

func (fi *DirNameOnlyFileInfo) Sys() any {
	return nil
}

func newDirNameOnlyFileInfo(name string, meta *FileMeta, fileOpener func() (afero.File, error)) FileMetaInfo {
	name = normalizeFilename(name)
	_, base := filepath.Split(name)

	m := meta.Copy()
	if m.Filename == "" {
		m.Filename = name
	}
	m.OpenFunc = fileOpener
	m.IsOrdered = false

	return NewFileMetaInfo(
		&DirNameOnlyFileInfo{name: base, modTime: htime.Now()},
		m,
	)
}

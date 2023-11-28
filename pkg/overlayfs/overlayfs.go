package overlayfs

import (
	"github.com/spf13/afero"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type AbsStatFss interface {
	// Abs returns an absolute path of file or dir.
	Abs(name string) []string

	Fss() []fs.StatFS
}

// OverlayFs is a filesystem that overlays multiple filesystems.
// It's by default a read-only filesystem.
// For all operations, the filesystems are checked in order until found.
type OverlayFs struct {
	fss []AbsStatFss
}

func (ofs *OverlayFs) Create(name string) (afero.File, error) {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Mkdir(name string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) MkdirAll(path string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) RemoveAll(path string) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Rename(oldname, newname string) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Name() string {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Chmod(name string, mode os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Chown(name string, uid, gid int) error {
	//TODO implement me
	panic("implement me")
}

func (ofs *OverlayFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	//TODO implement me
	panic("implement me")
}

// New creates a new OverlayFs with the given options.
func New(fss []AbsStatFss) *OverlayFs {
	return &OverlayFs{
		fss: fss,
	}
}

func (ofs *OverlayFs) Append(fss ...AbsStatFss) *OverlayFs {
	ofs.fss = append(ofs.fss, fss...)
	return ofs
}

// Open opens a file, returning it or an error, if any happens.
// If name is a directory, a *Dir is returned representing all directories matching name.
// Note that a *Dir must not be used after it's closed.
func (ofs *OverlayFs) Open(name string) (afero.File, error) {
	bfs, fi, err := ofs.stat(name)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, os.ErrInvalid
	}

	f, err := bfs.Open(name)
	if err != nil {
		return nil, err
	}

	return &File{File: f}, err
}

func (ofs *OverlayFs) Stat(name string) (os.FileInfo, error) {
	_, fi, err := ofs.stat(name)
	return fi, err
}

func (ofs *OverlayFs) stat(name string) (fs.StatFS, os.FileInfo, error) {
	for _, bfs := range ofs.fss {
		fss := bfs.Fss()
		for _, sfs := range fss {
			if fi, err := sfs.Stat(name); err == nil ||
				!os.IsNotExist(err) {
				return sfs, fi, nil
			}
		}
	}
	return nil, nil, os.ErrNotExist
}

func readFile(f fs.File) string {
	b, _ := io.ReadAll(f)
	return string(b)
}

func (ofs *OverlayFs) ReadDir(dirname string) ([]fs.FileInfo, error) {
	var fis []fs.FileInfo
	readDir := func(bfs AbsStatFss) error {
		dirs := bfs.Abs(dirname)
		for _, dir := range dirs {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			fis = merge(fis, files)
		}

		return nil
	}

	for _, bfs := range ofs.fss {
		if err := readDir(bfs); err != nil {
			return nil, err
		}
	}

	sort.Slice(fis, func(i, j int) bool {
		return fis[i].Name() < fis[j].Name()
	})

	return fis, nil
}

func merge(upper, lower []fs.FileInfo) []fs.FileInfo {
	for _, lfi := range lower {
		var found bool
		for _, ufi := range upper {
			if lfi.Name() == ufi.Name() {
				found = true
				break
			}
		}
		if !found {
			upper = append(upper, lfi)
		}
	}
	return upper
}

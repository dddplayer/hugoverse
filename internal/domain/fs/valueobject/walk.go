package valueobject

import (
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type (
	WalkFunc func(path string, info FileMetaInfo, err error) error
)

type Walkway struct {
	Fs       afero.Fs
	Root     string
	basePath string

	// May be pre-set
	fi         FileMetaInfo
	dirEntries []FileMetaInfo

	WalkFn WalkFunc

	// We may traverse symbolic links and bite ourself.
	Seen map[string]bool
}

func (w *Walkway) Walk() error {
	var fi FileMetaInfo
	if w.fi != nil {
		fi = w.fi
	} else {
		info, _, err := LstatIfPossible(w.Fs, w.Root)
		if err != nil {
			if os.IsNotExist(err) {
				return err
			}
			return w.WalkFn(w.Root, nil, fmt.Errorf("walk: %q: %w", w.Root, err))
		}
		fi = info.(FileMetaInfo)
	}

	if !fi.IsDir() {
		return w.WalkFn(w.Root, nil, errors.New("file to walk must be a directory"))
	}

	return w.walk(w.Root, fi, w.dirEntries, w.WalkFn)
}

// walk recursively descends path, calling walkFn.
// It follows symlinks if supported by the filesystem, but only the same path once.
func (w *Walkway) walk(path string, info FileMetaInfo, dirEntries []FileMetaInfo, walkFn WalkFunc) error {
	fmt.Println(">>> walk: ", path)
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return nil
	}

	meta := info.Meta()
	filename := meta.Filename

	if dirEntries == nil {
		f, err := w.Fs.Open(path)
		if err != nil {
			if w.checkErr(path, err) {
				return nil
			}
			return walkFn(path, info, fmt.Errorf("walk: open %q (%q): %w", path, w.Root, err))
		}

		fis, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			if w.checkErr(filename, err) {
				return nil
			}
			return walkFn(path, info, fmt.Errorf("walk: Readdir: %w", err))
		}

		fmt.Println(">>> dirEntries: ", fis, path)
		dirEntries = fileInfosToFileMetaInfos(fis)
	}

	// First add some metadata to the dir entries
	for _, fi := range dirEntries {
		fim := fi.(FileMetaInfo)

		meta := fim.Meta()

		// Note that we use the original TmplName even if it's a symlink.
		name := meta.Name
		if name == "" {
			name = fim.Name()
		}

		if name == "" {
			panic(fmt.Sprintf("[%s] no name set in %v", path, meta))
		}
		pathn := filepath.Join(path, name)

		pathMeta := pathn
		if w.basePath != "" {
			pathMeta = strings.TrimPrefix(pathn, w.basePath)
		}

		meta.Path = normalizeFilename(pathMeta)
		meta.PathWalk = pathn

		if fim.IsDir() && meta.IsSymlink && w.isSeen(meta.Filename) {
			// Prevent infinite recursion
			// Possible cyclic reference
			meta.SkipDir = true
		}
	}

	for _, fi := range dirEntries {
		fim := fi.(FileMetaInfo)
		meta := fim.Meta()

		if meta.SkipDir {
			continue
		}

		err := w.walk(meta.PathWalk, fim, nil, walkFn)
		if err != nil {
			if !fi.IsDir() || err != filepath.SkipDir {
				return err
			}
		}
	}

	return nil
}

func (w *Walkway) isSeen(filename string) bool {
	if filename == "" {
		return false
	}

	if w.Seen[filename] {
		return true
	}

	w.Seen[filename] = true
	return false
}

// checkErr returns true if the error is handled.
func (w *Walkway) checkErr(filename string, err error) bool {
	if os.IsNotExist(err) {
		// The file may be removed in process.
		// This may be a ERROR situation, but it is not possible
		// to determine as a general case.
		fmt.Printf("File %q not found, skipping.", filename)
		return true
	}

	return false
}

func fileInfosToFileMetaInfos(fis []os.FileInfo) []FileMetaInfo {
	fims := make([]FileMetaInfo, len(fis))
	for i, v := range fis {
		fims[i] = v.(FileMetaInfo)
	}
	return fims
}

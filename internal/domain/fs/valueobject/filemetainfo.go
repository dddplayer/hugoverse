package valueobject

import (
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/spf13/afero"
	"golang.org/x/text/unicode/norm"
	"os"
	"runtime"
	"strings"
)

type FileMetaInfo interface {
	os.FileInfo
	Meta() *FileMeta
}

func NewFileMetaInfo(fi os.FileInfo, m *FileMeta) FileMetaInfo {
	if m == nil {
		panic("FileMeta must be set")
	}
	if fim, ok := fi.(FileMetaInfo); ok {
		m.Merge(fim.Meta())
	}
	return &fileInfoMeta{FileInfo: fi, m: m}
}

func DecorateFileInfo(fi os.FileInfo, metaFs afero.Fs, opener func() (afero.File, error),
	filename, filepath string, inMeta *FileMeta) FileMetaInfo {
	var meta *FileMeta
	var fim FileMetaInfo

	filepath = strings.TrimPrefix(filepath, fs.FilepathSeparator)

	var ok bool
	if fim, ok = fi.(FileMetaInfo); ok {
		meta = fim.Meta()
	} else {
		meta = NewFileMeta()
		fim = NewFileMetaInfo(fi, meta)
	}

	if opener != nil {
		meta.OpenFunc = opener
	}
	if metaFs != nil {
		meta.Fs = metaFs
	}
	nfilepath := normalizeFilename(filepath)
	nfilename := normalizeFilename(filename)
	if nfilepath != "" {
		meta.Path = nfilepath
	}
	if nfilename != "" {
		meta.Filename = nfilename
	}

	meta.Merge(inMeta)

	return fim
}

func normalizeFilename(filename string) string {
	if filename == "" {
		return ""
	}
	if runtime.GOOS == "darwin" {
		// When a file system is HFS+, its filepath is in NFD form.
		return norm.NFC.String(filename)
	}
	return filename
}

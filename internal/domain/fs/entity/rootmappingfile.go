package entity

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	"github.com/spf13/afero"
	"os"
)

type rootMappingFile struct {
	afero.File
	fs   *RootMappingFs
	name string
	meta *valueobject.FileMeta
}

func (f *rootMappingFile) Close() error {
	if f.File == nil {
		return nil
	}
	return f.File.Close()
}

func (f *rootMappingFile) Name() string {
	return f.name
}

func (f *rootMappingFile) Readdir(count int) ([]os.FileInfo, error) {

	fmt.Println("Readdir: ", f.name)

	if f.File != nil {

		fis, err := f.File.Readdir(count)
		if err != nil {
			return nil, err
		}

		var result []os.FileInfo
		for _, fi := range fis {

			// TODO 8, understand how to make this work
			fim := valueobject.DecorateFileInfo(fi, f.fs, nil, "", "", f.meta)
			fmt.Println("Readdir fis : ", fim)

			if fim.IsDir() {
				result = append(result, fim)
			}
		}
		return result, nil
	}

	return f.fs.collectDirEntries(f.name)
}

func (f *rootMappingFile) Readdirnames(count int) ([]string, error) {
	dirs, err := f.Readdir(count)
	if err != nil {
		return nil, err
	}
	return fileInfosToNames(dirs), nil
}

func fileInfosToNames(fis []os.FileInfo) []string {
	names := make([]string, len(fis))
	for i, d := range fis {
		names[i] = d.Name()
	}
	return names
}

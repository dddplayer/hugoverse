package entity

import (
	"github.com/spf13/afero"
	"os"
)

type Fs struct {
	// Source is Hugo's source file system.
	// Note that this will always be a "plain" Afero filesystem:
	// * afero.OsFs when running in production
	// * afero.MemMapFs for many of the tests.
	source afero.Fs

	// PublishDir is where Hugo publishes its rendered content.
	// It's mounted inside publishDir (default /public).
	publishDir afero.Fs

	// WorkingDirReadOnly is a read-only file system
	// restricted to the project working dir.
	workingDirReadOnly afero.Fs
}

// NewFs creates a new Fs.
func NewFs(workingDir, publishDir string) *Fs {
	afs := afero.NewOsFs()
	workingFs := afero.NewBasePathFs(afs, workingDir)

	// Make sure we always have the /public folder ready to use.
	if err := workingFs.MkdirAll(publishDir, 0777); err != nil && !os.IsExist(err) {
		panic(err)
	}
	pubFs := afero.NewBasePathFs(workingFs, publishDir)

	return &Fs{
		source:             workingFs,
		publishDir:         pubFs,
		workingDirReadOnly: getWorkingDirFsReadOnly(workingFs, workingDir),
	}
}

func (f *Fs) Source() afero.Fs {
	return f.source
}

func (f *Fs) PublishDir() afero.Fs {
	return f.publishDir
}

func (f *Fs) WorkingDirReadOnly() afero.Fs {
	return f.workingDirReadOnly
}

func getWorkingDirFsReadOnly(base afero.Fs, workingDir string) afero.Fs {
	if workingDir == "" {
		return afero.NewReadOnlyFs(base)
	}
	return afero.NewBasePathFs(afero.NewReadOnlyFs(base), workingDir)
}

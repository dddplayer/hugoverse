package fs

import "github.com/spf13/afero"

// Fs holds the core filesystems used by Hugo.
type Fs interface {
	// Source is Hugo's source file system.
	// Note that this will always be a "plain" Afero filesystem:
	// * afero.OsFs when running in production
	// * afero.MemMapFs for many of the tests.
	Source() afero.Fs

	// PublishDir is where Hugo publishes its rendered content.
	// It's mounted inside publishDir (default /public).
	PublishDir() afero.Fs

	// WorkingDirReadOnly is a read-only file system
	// restricted to the project working dir.
	WorkingDirReadOnly() afero.Fs
}

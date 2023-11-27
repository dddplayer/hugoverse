package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	fsEntity "github.com/dddplayer/hugoverse/internal/domain/fs/entity"
)

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
	*Paths
	*fsEntity.BaseFs

	// The file systems to use
	Fs fs.Fs

	// The config provider to use
	Cfg config.Provider
}
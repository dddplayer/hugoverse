package entity

import (
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
)

type Paths struct {
	Fs         fs.Fs
	AllModules config.Modules
	WorkingDir string
}

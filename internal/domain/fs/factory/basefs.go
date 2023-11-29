package factory

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
	psEntity "github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
	"github.com/spf13/afero"
)

// NewBaseFS builds the filesystems used by Hugo given the paths and options provided.NewBase
func NewBaseFS(p *psEntity.Paths) (*entity.BaseFs, error) {
	fs := p.Fs

	publishFs := newBaseFileDecorator(fs.PublishDir())
	sourceFs := newBaseFileDecorator(afero.NewBasePathFs(fs.Source(), p.WorkingDir))

	b := &entity.BaseFs{
		SourceFs:  sourceFs,
		WorkDir:   fs.WorkingDirReadOnly(),
		PublishFs: publishFs,
	}

	builder := newSourceFilesystemsBuilder(p, b)
	sourceFilesystems, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("build filesystems: %w", err)
	}

	b.SourceFilesystems = sourceFilesystems
	b.TheBigFs = builder.theBigFs

	return b, nil
}

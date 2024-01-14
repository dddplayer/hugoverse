package factory

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/fs"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	psEntity "github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
	"github.com/dddplayer/hugoverse/pkg/overlayfs"
	"github.com/spf13/afero"
	"strings"
)

func newSourceFilesystem(name string, fs afero.Fs, dirs []valueobject.FileMetaInfo) *entity.SourceFilesystem {
	return &entity.SourceFilesystem{
		Name: name,
		Fs:   fs,
		Dirs: dirs,
	}
}

func newSourceFilesystemsBuilder(p *psEntity.Paths, b *entity.BaseFs) *sourceFilesystemsBuilder {
	sourceFs := NewBaseFileDecorator(p.Fs.Source())
	return &sourceFilesystemsBuilder{p: p, sourceFs: sourceFs, theBigFs: b.TheBigFs, result: &entity.SourceFilesystems{}}
}

type sourceFilesystemsBuilder struct {
	p        *psEntity.Paths
	sourceFs afero.Fs
	result   *entity.SourceFilesystems
	theBigFs *entity.FilesystemsCollector
}

func (b *sourceFilesystemsBuilder) Build() (*entity.SourceFilesystems, error) {
	if b.theBigFs == nil {
		// Modules - mounts <-> RootMappingFs - OverlayFS
		theBigFs, err := b.createMainOverlayFs(b.p)
		if err != nil {
			return nil, fmt.Errorf("create main fs: %w", err)
		}

		b.theBigFs = theBigFs
	}

	createView := func(componentID string, ofs *overlayfs.OverlayFs) *entity.SourceFilesystem {
		dirs := b.theBigFs.OverlayDirs[componentID]
		return newSourceFilesystem(componentID, afero.NewBasePathFs(ofs, componentID), dirs)
	}

	b.result.Layouts = createView(fs.ComponentFolderLayouts, b.theBigFs.OverlayMounts)
	b.result.Content = createView(fs.ComponentFolderContent, b.theBigFs.OverlayMountsContent)

	return b.result, nil
}

func (b *sourceFilesystemsBuilder) createMainOverlayFs(p *psEntity.Paths) (*entity.FilesystemsCollector, error) {
	collector := &entity.FilesystemsCollector{
		SourceProject:        b.sourceFs,
		OverlayMounts:        overlayfs.New([]overlayfs.AbsStatFss{}),
		OverlayMountsContent: overlayfs.New([]overlayfs.AbsStatFss{}),
		OverlayDirs:          make(map[string][]valueobject.FileMetaInfo),
	}
	err := b.createOverlayFs(collector, p)

	return collector, err
}

func (b *sourceFilesystemsBuilder) createOverlayFs(collector *entity.FilesystemsCollector, path *psEntity.Paths) error {
	for _, md := range path.AllModules {

		var (
			fromTo        []valueobject.RootMapping
			fromToContent []valueobject.RootMapping
		)

		for _, mount := range md.Mounts() {
			rm := valueobject.RootMapping{
				From: mount.Target, // content
				To:   mount.Source, // mycontent
				Meta: &valueobject.FileMeta{
					Classifier: valueobject.ContentClassContent,
				},
			}

			isContentMount := b.isContentMount(mount.Target)
			if isContentMount {
				fromToContent = append(fromToContent, rm)
			} else if b.isStaticMount(mount.Target) {
				continue
			} else {
				fromTo = append(fromTo, rm)
			}
		}

		rmfs := newRootMappingFs(collector.SourceProject, fromTo...)
		rmfsContent := newRootMappingFs(collector.SourceProject, fromToContent...)

		collector.AddDirs(rmfs)        // add other folders, /layouts etc
		collector.AddDirs(rmfsContent) // only has /content, why need to go through all components?

		collector.OverlayMounts = collector.OverlayMounts.Append(rmfs)
		collector.OverlayMountsContent = collector.OverlayMountsContent.Append(rmfsContent)
	}

	return nil
}

func (b *sourceFilesystemsBuilder) isContentMount(target string) bool {
	return strings.HasPrefix(target, fs.ComponentFolderContent)
}

func (b *sourceFilesystemsBuilder) isStaticMount(target string) bool {
	return strings.HasPrefix(target, fs.ComponentFolderStatic)
}

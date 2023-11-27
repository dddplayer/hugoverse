package factory

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/fs/entity"
	"github.com/dddplayer/hugoverse/internal/domain/fs/valueobject"
	psEntity "github.com/dddplayer/hugoverse/internal/domain/pathspec/entity"
	"github.com/dddplayer/hugoverse/pkg/overlayfs"
	"github.com/spf13/afero"
)

func newSourceFilesystemsBuilder(p *psEntity.Paths, b *entity.BaseFs) *sourceFilesystemsBuilder {
	sourceFs := newBaseFileDecorator(p.Fs.Source())
	return &sourceFilesystemsBuilder{p: p, sourceFs: sourceFs, theBigFs: b.TheBigFs, result: &entity.SourceFilesystems{}}
}

type sourceFilesystemsBuilder struct {
	p        *psEntity.Paths
	sourceFs afero.Fs
	result   *entity.SourceFilesystems
	theBigFs *valueobject.FilesystemsCollector
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
	//
	//createView := func(componentID string) *SourceFilesystem {
	//	if b.theBigFs == nil || b.theBigFs.overlayMounts == nil {
	//		return b.newSourceFilesystem(componentID, hugofs.NoOpFs, nil)
	//	}
	//
	//	dirs := b.theBigFs.overlayDirs[componentID]
	//
	//	return b.newSourceFilesystem(componentID, afero.NewBasePathFs(b.theBigFs.overlayMounts, componentID), dirs)
	//}
	//
	//b.result.Archetypes = createView(files.ComponentFolderArchetypes)
	//b.result.Layouts = createView(files.ComponentFolderLayouts)
	//b.result.Assets = createView(files.ComponentFolderAssets)
	//b.result.ResourcesCache = b.theBigFs.overlayResources
	//
	//// Data, i18n and content cannot use the overlay fs
	//dataDirs := b.theBigFs.overlayDirs[files.ComponentFolderData]
	//dataFs, err := hugofs.NewSliceFs(dataDirs...)
	//if err != nil {
	//	return nil, err
	//}
	//
	//b.result.Data = b.newSourceFilesystem(files.ComponentFolderData, dataFs, dataDirs)
	//
	//i18nDirs := b.theBigFs.overlayDirs[files.ComponentFolderI18n]
	//i18nFs, err := hugofs.NewSliceFs(i18nDirs...)
	//if err != nil {
	//	return nil, err
	//}
	//b.result.I18n = b.newSourceFilesystem(files.ComponentFolderI18n, i18nFs, i18nDirs)
	//
	//contentDirs := b.theBigFs.overlayDirs[files.ComponentFolderContent]
	//contentBfs := afero.NewBasePathFs(b.theBigFs.overlayMountsContent, files.ComponentFolderContent)
	//
	//contentFs, err := hugofs.NewLanguageFs(b.p.LanguagesDefaultFirst.AsOrdinalSet(), contentBfs)
	//if err != nil {
	//	return nil, fmt.Errorf("create content filesystem: %w", err)
	//}
	//
	//b.result.Content = b.newSourceFilesystem(files.ComponentFolderContent, contentFs, contentDirs)
	//
	//b.result.Work = afero.NewReadOnlyFs(b.theBigFs.overlayFull)
	//
	//// Create static filesystem(s)
	//ms := make(map[string]*SourceFilesystem)
	//b.result.Static = ms
	//b.result.StaticDirs = b.theBigFs.overlayDirs[files.ComponentFolderStatic]
	//
	//bfs := afero.NewBasePathFs(b.theBigFs.overlayMountsStatic, files.ComponentFolderStatic)
	//ms[""] = b.newSourceFilesystem(files.ComponentFolderStatic, bfs, b.result.StaticDirs)

	return b.result, nil
}

func (b *sourceFilesystemsBuilder) createMainOverlayFs(p *psEntity.Paths) (*valueobject.FilesystemsCollector, error) {
	collector := &valueobject.FilesystemsCollector{
		OverlayMountsContent: overlayfs.New([]overlayfs.AbsStatFss{}),
	}
	err := b.createOverlayFs(collector, p)

	return collector, err
}

func (b *sourceFilesystemsBuilder) createOverlayFs(collector *valueobject.FilesystemsCollector, path *psEntity.Paths) error {
	for _, md := range path.AllModules {
		var fromToContent []valueobject.RootMapping

		for _, mount := range md.Mounts() {
			rm := valueobject.RootMapping{
				From: mount.Target, // content
				To:   mount.Source, // mycontent
			}
			fromToContent = append(fromToContent, rm)
		}
		rmfsContent := newRootMappingFs(fromToContent...)
		collector.OverlayMountsContent = collector.OverlayMountsContent.Append(rmfsContent)
	}
	return nil
}

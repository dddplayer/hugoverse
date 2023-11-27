package valueobject

import "github.com/dddplayer/hugoverse/pkg/overlayfs"

type FilesystemsCollector struct {
	OverlayMountsContent *overlayfs.OverlayFs
}

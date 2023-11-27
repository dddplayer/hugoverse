package valueobject

import "github.com/dddplayer/hugoverse/internal/domain/config"

// moduleAdapter implemented Module interface
type moduleAdapter struct {
	projectMod bool
	owner      config.Module
	mounts     []config.Mount
	config     config.ModuleConfig
}

func (m *moduleAdapter) Config() config.ModuleConfig {
	return m.config
}
func (m *moduleAdapter) Mounts() []config.Mount {
	return m.mounts
}
func (m *moduleAdapter) Owner() config.Module {
	return m.owner
}

// Module folder structure
const (
	ComponentFolderArchetypes = "archetypes"
	ComponentFolderStatic     = "static"
	ComponentFolderLayouts    = "layouts"
	ComponentFolderContent    = "content"
	ComponentFolderData       = "data"
	ComponentFolderAssets     = "assets"
	ComponentFolderI18n       = "i18n"
)

// ApplyProjectConfigDefaults applies default/missing module configuration for
// the main project.
func ApplyProjectConfigDefaults(mod config.Module) {
	projectMod := mod.(*moduleAdapter)

	type dirKeyComponent struct {
		key          string
		component    string
		multilingual bool
	}

	dirKeys := []dirKeyComponent{
		{"contentDir", ComponentFolderContent, true},
		{"dataDir", ComponentFolderData, false},
		{"layoutDir", ComponentFolderLayouts, false},
		{"i18nDir", ComponentFolderI18n, false},
		{"archetypeDir", ComponentFolderArchetypes,
			false},
		{"assetDir", ComponentFolderAssets, false},
		{"", ComponentFolderStatic, false},
	}

	var mounts []config.Mount
	for _, d := range dirKeys {
		if d.multilingual {
			// based on language content configuration
			// multiple language has multiple source folders
			if d.component == ComponentFolderContent {
				mounts = append(mounts, config.Mount{
					Lang:   "en",
					Source: "mycontent",
					Target: d.component},
				)
			}
		} else {
			mounts = append(mounts,
				config.Mount{
					Source: d.component,
					Target: d.component},
			)
		}
	}

	projectMod.mounts = mounts
}

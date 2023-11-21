package valueobject

type Mount struct {
	// relative path in source repo, e.g. "scss"
	Source string
	// relative target path, e.g. "assets/bootstrap/scss"
	Target string
	// any language code associated with this mount.
	Lang string
}

type Import struct {
	// Module path
	Path string
}

// ModuleConfig holds a module config.
type ModuleConfig struct {
	Mounts  []Mount
	Imports []Import
}

type Module interface {
	// Config The decoded module config and mounts.
	Config() ModuleConfig
	// Owner In the dependency tree, this is the first module that defines this module
	// as a dependency.
	Owner() Module
	// Mounts Any directory remappings.
	Mounts() []Mount
}

type Modules []Module

// moduleAdapter implemented Module interface
type moduleAdapter struct {
	projectMod bool
	owner      Module
	mounts     []Mount
	config     ModuleConfig
}

func (m *moduleAdapter) Config() ModuleConfig {
	return m.config
}
func (m *moduleAdapter) Mounts() []Mount {
	return m.mounts
}
func (m *moduleAdapter) Owner() Module {
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
func ApplyProjectConfigDefaults(mod Module) {
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

	var mounts []Mount
	for _, d := range dirKeys {
		if d.multilingual {
			// based on language content configuration
			// multiple language has multiple source folders
			if d.component == ComponentFolderContent {
				mounts = append(mounts, Mount{
					Lang:   "en",
					Source: "mycontent",
					Target: d.component})
			}
		} else {
			mounts = append(mounts,
				Mount{
					Source: d.component,
					Target: d.component})
		}
	}

	projectMod.mounts = mounts
}

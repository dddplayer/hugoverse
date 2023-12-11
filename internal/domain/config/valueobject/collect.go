package valueobject

import (
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/config"
)

type ModuleCollector struct {
	Modules config.Modules
}

func NewModuleCollector() *ModuleCollector {
	return &ModuleCollector{
		Modules: config.Modules{},
	}
}

func (mc *ModuleCollector) CollectModules(modConfig config.ModuleConfig, hookBeforeFinalize func(m config.Modules)) {
	projectMod := &moduleAdapter{
		projectMod: true,
		config:     modConfig,
	}

	// module structure, [project, others...]
	mc.addAndRecurse(projectMod)

	// Add the project mod on top.
	mc.Modules = append(config.Modules{projectMod}, mc.Modules...)

	if hookBeforeFinalize != nil {
		hookBeforeFinalize(mc.Modules)
	}
}

// addAndRecurse Project Imports -> Import imports
func (mc *ModuleCollector) addAndRecurse(owner *moduleAdapter) {
	moduleConfig := owner.Config()

	// theme may depend on other theme
	for _, moduleImport := range moduleConfig.Imports {
		tc := mc.add(owner, moduleImport)
		if tc == nil {
			continue
		}
		// tc is mytheme with no config file
		mc.addAndRecurse(tc)
	}
}

func (mc *ModuleCollector) add(owner *moduleAdapter, moduleImport config.Import) *moduleAdapter {
	fmt.Printf("--- start to create `%s` module\n", moduleImport.Path)
	ma := &moduleAdapter{
		owner: owner,
		// In the example, "mytheme" has no other import
		// In the real world, we need to parse the theme config and download the theme repo
		config: config.ModuleConfig{},
	}
	mc.Modules = append(mc.Modules, ma)
	return ma
}

package valueobject

import "fmt"

type ModuleCollector struct {
	Modules Modules
}

func NewModuleCollector() *ModuleCollector {
	return &ModuleCollector{
		Modules: Modules{},
	}
}

func (mc *ModuleCollector) CollectModules(modConfig ModuleConfig, hookBeforeFinalize func(m Modules)) {
	projectMod := &moduleAdapter{
		projectMod: true,
		config:     modConfig,
	}

	// module structure, [project, others...]
	mc.addAndRecurse(projectMod)

	// Add the project mod on top.
	mc.Modules = append(Modules{projectMod}, mc.Modules...)

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

func (mc *ModuleCollector) add(owner *moduleAdapter,
	moduleImport Import) *moduleAdapter {

	fmt.Printf("start to create `%s` module\n",
		moduleImport.Path)
	ma := &moduleAdapter{
		owner: owner,
		// in the example, mytheme has no other import
		config: ModuleConfig{},
	}
	mc.Modules = append(mc.Modules, ma)
	return ma
}

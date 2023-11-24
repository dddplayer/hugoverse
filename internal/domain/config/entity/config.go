package entity

import (
	"bytes"
	"fmt"
	"github.com/dddplayer/hugoverse/internal/domain/config"
	"github.com/dddplayer/hugoverse/internal/domain/config/factory"
	"github.com/dddplayer/hugoverse/internal/domain/config/valueobject"
	"github.com/dddplayer/hugoverse/pkg/log"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Format 文件格式类型
type Format string

// TOML 支持的格式，为简单示例，只支持TOML格式
const (
	TOML Format = "toml"
)

type Config struct {
	Path    string
	Modules valueobject.Modules
	Logger  log.Logger
}

func (c *Config) Load() (config.Provider, error) {
	m, err := c.loadConfigFromDisk()
	if err != nil {
		return nil, err
	}

	provider := factory.New()
	provider.Set("", m)
	provider.SetDefaults(config.Params{
		"Path":    c.Path,
		"timeout": "30s",
	})

	theme := provider.GetString("theme")
	if theme != "" {
		modules, err := c.loadModules(theme)
		if err != nil {
			return nil, err
		}
		provider.Set("modules", modules)
		for _, m := range modules {
			c.Logger.Printf("module: %+v", m)
		}
	}

	return provider, nil
}

func (c *Config) loadModules(theme string) (valueobject.Modules, error) {
	// project module config
	projModuleConfig := valueobject.ModuleConfig{}
	imports := []string{theme}
	for _, imp := range imports {
		projModuleConfig.Imports = append(
			projModuleConfig.Imports, valueobject.Import{
				Path: imp,
			})
	}

	mc := valueobject.NewModuleCollector()
	// Need to run these after the modules are loaded, but before
	// they are finalized.
	collectHook := func(mods valueobject.Modules) {
		// Apply default project mounts.
		// Default folder structure for hugo project
		valueobject.ApplyProjectConfigDefaults(mods[0])
	}
	mc.CollectModules(projModuleConfig, collectHook)

	return mc.Modules, nil
}

func (c *Config) loadConfigFromDisk() (map[string]any, error) {
	content, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}

	configData := bytes.TrimSuffix(content, []byte("\n"))
	format := FormatFromString(path.Base(c.Path))
	m := make(map[string]any)

	if err := UnmarshalTo(configData, format, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
// Hugo 实现
func FormatFromString(formatStr string) Format {
	formatStr = strings.ToLower(formatStr)
	if strings.Contains(formatStr, ".") {
		// Assume a filename
		formatStr = strings.TrimPrefix(
			filepath.Ext(formatStr), ".")
	}
	switch formatStr {
	case "toml":
		return TOML
	}

	return ""
}

// UnmarshalTo unmarshals data in format f into v.
func UnmarshalTo(data []byte, f Format, v any) error {
	var err error

	switch f {
	case TOML:
		err = toml.Unmarshal(data, v)

	default:
		return fmt.Errorf(
			"unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return err
}

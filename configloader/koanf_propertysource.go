package configloader

import (
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"strings"
)

const (
	propertyFileName = "application.yaml"
	propertyFilePath = "PROPERTY_FILE_PATH"
)

type PropertySource struct {
	// Provider represents a configuration provider. Providers can
	// read configuration from a source (file, HTTP etc.)
	// Methods Read and ReadBytes must return unflatten map or string keys should not be flat delimited
	Provider PropertyProvider

	// Parser represents a configuration format parser.
	Parser koanf.Parser
}

// YamlPropertySourceParams allows to configure YamlPropertySource.
type YamlPropertySourceParams struct {
	//ConfigFilePath to determine path to application.yaml file
	ConfigFilePath string
}

func YamlPropertySource(params ...YamlPropertySourceParams) *PropertySource {
	configFilePath := getConfigFilePath(params)
	return &PropertySource{
		Provider: AsPropertyProvider(file.Provider(configFilePath)),
		Parser:   flattenYAMLParser{},
	}
}

func getConfigFilePath(params []YamlPropertySourceParams) string {
	path := ""
	if params != nil && params[0].ConfigFilePath != "" {
		path = params[0].ConfigFilePath
	} else {
		pathToFolder := os.Getenv(propertyFilePath)
		path = pathToFolder + propertyFileName
	}
	return path
}

func EnvPropertySource() *PropertySource {
	return &PropertySource{Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil}
}

func BasePropertySources(params ...YamlPropertySourceParams) []*PropertySource {
	return []*PropertySource{YamlPropertySource(params...), EnvPropertySource()}
}

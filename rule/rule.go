package rule

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Config []ConfigItem `yaml:"config"`
}

type ConfigItem struct {
	Name string `yaml:"name"`
}

func Add(yamlContent string, rule string) (string, error) {
	var config Config
	commentMap := yaml.CommentMap{}

	decoder := yaml.NewDecoder(strings.NewReader(yamlContent), yaml.CommentToMap(commentMap))
	if err := decoder.Decode(&config); err != nil {
		return "", fmt.Errorf("failed to parse YAML: %w", err)
	}

	config.Config = append(config.Config, ConfigItem{
		Name: rule,
	})

	updatedYAML, err := yaml.MarshalWithOptions(
		&config,
		yaml.Indent(2),
		yaml.IndentSequence(true),
		yaml.WithComment(commentMap),
	)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return string(updatedYAML), nil
}

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

func GetAddedRules(old, new string) ([]string, error) {
	var oldConfig, newConfig Config
	if err := yaml.Unmarshal([]byte(old), &oldConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal old YAML: %w", err)
	}
	if err := yaml.Unmarshal([]byte(new), &newConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal new YAML: %w", err)
	}

	oldRules := make(map[string]struct{})
	for _, item := range oldConfig.Config {
		oldRules[item.Name] = struct{}{}
	}

	addedRules := []string{}
	for _, item := range newConfig.Config {
		if _, exists := oldRules[item.Name]; !exists {
			addedRules = append(addedRules, item.Name)
		}
	}

	return addedRules, nil
}

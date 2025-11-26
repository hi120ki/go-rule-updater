package env

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Environment             Environment `envconfig:"ENVIRONMENT" default:"development"`
	Owner                   string      `envconfig:"OWNER" default:"hi120ki"`
	Repository              string      `envconfig:"REPOSITORY" default:"go-rule-updater"`
	RulePath                string      `envconfig:"RULE_PATH" default:"rule.yaml"`
	BaseBranch              string      `envconfig:"BASE_BRANCH" default:"main"`
	GitHubToken             string      `envconfig:"GITHUB_TOKEN"`
	GitHubAppID             int64       `envconfig:"GITHUB_APP_ID" `
	GitHubAppInstallationID int64       `envconfig:"GITHUB_APP_INSTALLATION_ID" `
	GitHubAppPrivateKey     string      `envconfig:"GITHUB_APP_PRIVATE_KEY"`
}

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

func Load() (*Env, error) {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}
	return &env, nil
}

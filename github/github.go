package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v79/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/hi120ki/go-rule-updater/env"
)

type GitHub struct {
	client   *github.Client
	clientV4 *githubv4.Client
}

func NewClient(ctx context.Context, cfg *env.Env) (*GitHub, error) {
	if cfg.Environment == env.EnvironmentProduction {
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, cfg.GitHubAppID, cfg.GitHubAppInstallationID, cfg.GitHubAppPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create ghinstallation transport: %w", err)
		}
		itrClient := &http.Client{Transport: itr, Timeout: 5 * time.Second}
		return &GitHub{
			client:   github.NewClient(itrClient),
			clientV4: githubv4.NewClient(itrClient),
		}, nil
	}

	if cfg.GitHubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is required")
	}
	itrClient := oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GitHubToken},
		),
	)
	return &GitHub{
		client:   github.NewClient(itrClient),
		clientV4: githubv4.NewClient(itrClient),
	}, nil
}

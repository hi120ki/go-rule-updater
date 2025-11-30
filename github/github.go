package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v79/github"
	"github.com/octo-sts/app/pkg/gcpkms"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/hi120ki/go-rule-updater/env"
)

type GitHub struct {
	client   *github.Client
	clientV4 *githubv4.Client
}

func NewClient(ctx context.Context, cfg *env.Env) (*GitHub, error) {
	if cfg.GitHubAppKMSKeyPath != "" {
		kmsClient, err := kms.NewKeyManagementClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create kms client: %w", err)
		}

		signer, err := gcpkms.New(ctx, kmsClient, cfg.GitHubAppKMSKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create kms signer: %w", err)
		}

		atr, err := ghinstallation.NewAppsTransportWithOptions(http.DefaultTransport, cfg.GitHubAppID, ghinstallation.WithSigner(signer))
		if err != nil {
			return nil, fmt.Errorf("failed to create ghinstallation transport: %w", err)
		}

		itr := ghinstallation.NewFromAppsTransport(atr, cfg.GitHubAppInstallationID)

		itrClient := &http.Client{Transport: itr, Timeout: 5 * time.Second}
		return &GitHub{
			client:   github.NewClient(itrClient),
			clientV4: githubv4.NewClient(itrClient),
		}, nil
	}

	if cfg.GitHubAppPrivateKey != "" {
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

	if cfg.GitHubToken != "" {
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

	return nil, fmt.Errorf("no GitHub authentication method provided")
}

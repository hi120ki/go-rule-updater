package github_test

import (
	"context"
	"os"
	"testing"

	"github.com/hi120ki/go-rule-updater/env"
	"github.com/hi120ki/go-rule-updater/github"
)

func TestGitHub_GetLatestCommitSHA(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg *env.Env
		// Named input parameters for target function.
		owner   string
		repo    string
		branch  string
		wantErr bool
	}{
		{
			name: "Get latest commit SHA of existing branch",
			cfg: &env.Env{
				Environment: env.EnvironmentDevelopment,
				GitHubToken: os.Getenv("GITHUB_TOKEN"),
			},
			owner:   "hi120ki",
			repo:    "go-rule-updater",
			branch:  "main",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := github.NewClient(context.Background(), tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, gotErr := g.GetLatestCommitSHA(context.Background(), tt.owner, tt.repo, tt.branch)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetLatestCommitSHA() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetLatestCommitSHA() succeeded unexpectedly")
			}
		})
	}
}

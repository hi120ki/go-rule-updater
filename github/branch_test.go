package github_test

import (
	"context"
	"testing"

	"github.com/hi120ki/go-rule-updater/env"
	"github.com/hi120ki/go-rule-updater/github"
)

func TestGitHub_CreateBranch(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cfg *env.Env
		// Named input parameters for target function.
		owner      string
		repo       string
		branchName string
		baseBranch string
		wantErr    bool
	}{
		/*
			{
				name:        "Create branch in development environment",
				cfg: &env.Env{
					Environment: env.EnvironmentDevelopment,
				},
				owner:       "hi120ki",
				repo:        "go-rule-updater",
				branchName:  "test-branch",
				baseBranch:  "main",
				wantErr:     false,
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := github.NewClient(context.Background(), tt.cfg)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := g.CreateBranch(context.Background(), tt.owner, tt.repo, tt.branchName, tt.baseBranch)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateBranch() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateBranch() succeeded unexpectedly")
			}
		})
	}
}

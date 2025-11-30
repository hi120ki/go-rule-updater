package github_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v79/github"
	"github.com/hi120ki/go-rule-updater/env"
	ghclient "github.com/hi120ki/go-rule-updater/github"
)

func setupTestClient(t *testing.T) *ghclient.GitHub {
	t.Helper()
	cfg := &env.Env{
		Environment: env.EnvironmentDevelopment,
	}
	gh, err := ghclient.NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Failed to create GitHub client: %v", err)
	}
	return gh
}

func TestIsConflicting(t *testing.T) {
	tests := []struct {
		name     string
		pr       *github.PullRequest
		expected bool
	}{
		{
			name: "conflicting PR",
			pr: &github.PullRequest{
				Mergeable:      github.Ptr(false),
				MergeableState: github.Ptr("dirty"),
			},
			expected: true,
		},
		{
			name: "clean PR",
			pr: &github.PullRequest{
				Mergeable:      github.Ptr(true),
				MergeableState: github.Ptr("clean"),
			},
			expected: false,
		},
		{
			name: "unstable PR",
			pr: &github.PullRequest{
				Mergeable:      github.Ptr(false),
				MergeableState: github.Ptr("unstable"),
			},
			expected: false,
		},
		{
			name: "blocked PR",
			pr: &github.PullRequest{
				Mergeable:      github.Ptr(false),
				MergeableState: github.Ptr("blocked"),
			},
			expected: false,
		},
	}

	gh := &ghclient.GitHub{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gh.IsConflicting(tt.pr)
			if result != tt.expected {
				t.Errorf("IsConflicting() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListOpenPullRequests(t *testing.T) {
	ctx := context.Background()
	gh := setupTestClient(t)

	prs, err := gh.ListOpenPullRequests(ctx, "hi120ki", "go-rule-updater")
	if err != nil {
		t.Fatalf("ListOpenPullRequests() error = %v", err)
	}

	t.Logf("Found %d open PRs", len(prs))
}

func TestGetPullRequest(t *testing.T) {
	ctx := context.Background()
	gh := setupTestClient(t)

	pr, err := gh.GetPullRequest(ctx, "hi120ki", "go-rule-updater", 2)
	if err != nil {
		t.Fatalf("GetPullRequest() error = %v", err)
	}

	t.Logf("PR #%d: %s (mergeable: %v, state: %s)",
		pr.GetNumber(),
		pr.GetTitle(),
		pr.GetMergeable(),
		pr.GetMergeableState())
}

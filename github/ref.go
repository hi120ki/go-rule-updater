package github

import (
	"context"
	"fmt"
)

func (g *GitHub) GetLatestCommitSHA(ctx context.Context, owner, repo, branch string) (string, error) {
	ref, _, err := g.client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		return "", fmt.Errorf("failed to get ref for branch %s: %w", branch, err)
	}

	if ref == nil || ref.Object == nil {
		return "", fmt.Errorf("branch %s not found", branch)
	}

	sha := ref.Object.GetSHA()
	if sha == "" {
		return "", fmt.Errorf("branch %s has empty SHA", branch)
	}

	return sha, nil
}

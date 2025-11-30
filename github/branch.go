package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v79/github"
)

func (g *GitHub) CreateBranch(ctx context.Context, owner, repo, branchName, baseBranch string) error {
	ref, _, err := g.client.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch %s: %w", baseBranch, err)
	}

	baseSHA := ref.Object.GetSHA()
	if baseSHA == "" {
		return fmt.Errorf("base branch %s has no SHA", baseBranch)
	}

	_, _, err = g.client.Git.CreateRef(ctx, owner, repo, github.CreateRef{
		Ref: "refs/heads/" + branchName,
		SHA: baseSHA,
	})
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %w", branchName, err)
	}

	return nil
}

// UpdateBranchRef updates a branch reference to point to a new SHA (force update)
func (g *GitHub) UpdateBranchRef(ctx context.Context, owner, repo, branchName, sha string, force bool) error {
	ref := "heads/" + branchName
	updateRef := github.UpdateRef{
		SHA:   sha,
		Force: &force,
	}
	_, _, err := g.client.Git.UpdateRef(ctx, owner, repo, ref, updateRef)
	if err != nil {
		return fmt.Errorf("failed to update branch %s to SHA %s: %w", branchName, sha, err)
	}
	return nil
}

package github

import (
	"context"

	"github.com/google/go-github/v79/github"
)

func (g *GitHub) CreatePullRequest(ctx context.Context, owner, repo, title, head, base, body string) (*github.PullRequest, error) {
	pr, _, err := g.client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: github.Ptr(title),
		Head:  github.Ptr(head),
		Base:  github.Ptr(base),
		Body:  github.Ptr(body),
	})
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (g *GitHub) CreatePullRequestComment(ctx context.Context, owner, repo string, number int, message string) error {
	_, _, err := g.client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
		Body: github.Ptr(message),
	})
	return err
}

func (g *GitHub) MergePullRequest(ctx context.Context, owner, repo string, number int) error {
	_, _, err := g.client.PullRequests.Merge(ctx, owner, repo, number, "", &github.PullRequestOptions{})
	return err
}

// ListOpenPullRequests returns all open pull requests for the repository
func (g *GitHub) ListOpenPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	prs, _, err := g.client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
	})
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// GetPullRequest retrieves a specific pull request
func (g *GitHub) GetPullRequest(ctx context.Context, owner, repo string, number int) (*github.PullRequest, error) {
	pr, _, err := g.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

// UpdatePullRequestBranch updates the PR branch with the latest base branch changes
func (g *GitHub) UpdatePullRequestBranch(ctx context.Context, owner, repo string, number int) error {
	_, _, err := g.client.PullRequests.UpdateBranch(ctx, owner, repo, number, &github.PullRequestBranchUpdateOptions{})
	return err
}

// IsConflicting checks if a PR has merge conflicts
func (g *GitHub) IsConflicting(pr *github.PullRequest) bool {
	return !pr.GetMergeable() && pr.GetMergeableState() == "dirty"
}

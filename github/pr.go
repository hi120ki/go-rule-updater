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

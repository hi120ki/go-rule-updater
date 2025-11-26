package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v79/github"
)

func (g *GitHub) GetFile(ctx context.Context, owner, repo, path, ref string) (string, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	fileContent, _, _, err := g.client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return "", fmt.Errorf("failed to get file %s at ref %s: %w", path, ref, err)
	}

	if fileContent == nil {
		return "", fmt.Errorf("file %s not found at ref %s", path, ref)
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return content, nil
}
